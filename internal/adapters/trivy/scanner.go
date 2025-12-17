package trivy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/subilmondesir/podman-panoptic/internal/core/domain"
	"github.com/subilmondesir/podman-panoptic/internal/core/ports"
)

// Scanner implémente l'interface VulnerabilityScanner via le binaire Trivy CLI
type Scanner struct {
	binaryPath string
}

// NewScanner tente de trouver Trivy dans le PATH
func NewScanner() *Scanner {
	path, err := exec.LookPath("trivy")
	if err != nil {
		return nil // Trivy n'est pas installé
	}
	return &Scanner{binaryPath: path}
}

// IsAvailable retourne true si Trivy est installé
func (s *Scanner) IsAvailable() bool {
	return s != nil && s.binaryPath != ""
}

// ScanImage exécute trivy image avec une stratégie de fallback intelligente
func (s *Scanner) ScanImage(ctx context.Context, imageName string) ([]domain.Vulnerability, error) {
	if !s.IsAvailable() {
		return nil, fmt.Errorf("trivy non disponible")
	}

	// Nettoyage préventif du nom de l'image (ex: "localhost/alpine" -> "alpine")
	cleanImageName := strings.TrimPrefix(imageName, "localhost/")

	// Stratégie 1: Scan Standard (Trivy auto-détection)
	vulns, err := s.runTrivyCommand(ctx, cleanImageName, false)
	if err == nil {
		return vulns, nil
	}

	// Stratégie 2: Fallback Podman Explicite (Si le scan standard échoue)
	// Utile pour les configurations Rootless ou custom sockets
	vulnsFallback, errFallback := s.runTrivyCommand(ctx, cleanImageName, true)
	if errFallback == nil {
		return vulnsFallback, nil
	}

	// Si tout échoue, on retourne l'erreur initiale
	return nil, fmt.Errorf("scan échoué: %w", err)
}

// runTrivyCommand exécute la commande avec ou sans le flag podman
func (s *Scanner) runTrivyCommand(ctx context.Context, image string, usePodmanSrc bool) ([]domain.Vulnerability, error) {
	args := []string{
		"image",
		"--format", "json",
		"--quiet",
		"--severity", "HIGH,CRITICAL",
		"--scanners", "vuln",
	}

	if usePodmanSrc {
		args = append(args, "--image-src", "podman")
	}

	args = append(args, image)

	cmd := exec.CommandContext(ctx, s.binaryPath, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("err: %v | stderr: %s", err, stderr.String())
	}

	return parseTrivyOutput(stdout.Bytes())
}

// Structures internes pour mapper le JSON de Trivy
type trivyReport struct {
	Results []struct {
		Target          string `json:"Target"`
		Vulnerabilities []struct {
			VulnerabilityID  string `json:"VulnerabilityID"`
			PkgName          string `json:"PkgName"`
			InstalledVersion string `json:"InstalledVersion"`
			FixedVersion     string `json:"FixedVersion"`
			Title            string `json:"Title"`
			Description      string `json:"Description"`
			Severity         string `json:"Severity"`
			PrimaryURL       string `json:"PrimaryURL"`
		} `json:"Vulnerabilities"`
	} `json:"Results"`
}

func parseTrivyOutput(data []byte) ([]domain.Vulnerability, error) {
	// Cas limite : Sortie vide
	if len(data) == 0 {
		return []domain.Vulnerability{}, nil
	}

	var report trivyReport
	if err := json.Unmarshal(data, &report); err != nil {
		// Parfois Trivy renvoie du texte simple en cas d'erreur logique non fatale
		return nil, fmt.Errorf("parsing trivy json: %w", err)
	}

	var vulns []domain.Vulnerability
	for _, result := range report.Results {
		for _, v := range result.Vulnerabilities {
			vulns = append(vulns, domain.Vulnerability{
				ID:          v.VulnerabilityID,
				Severity:    domain.Severity(v.Severity),
				Title:       v.Title,
				Description: v.Description,
				Package:     v.PkgName,
				Version:     v.InstalledVersion,
				FixedIn:     v.FixedVersion,
				References:  []string{v.PrimaryURL},
			})
		}
	}
	return vulns, nil
}

// Vérification interface
var _ ports.VulnerabilityScanner = (*Scanner)(nil)
