package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/subilmondesir/podman-panoptic/internal/core/domain"
	"github.com/subilmondesir/podman-panoptic/internal/core/ports"
)

// ProgressCallback est une fonction appelée à chaque étape du scan (pour l'UI)
type ProgressCallback func(current, total int, message string)

// AuditService orchestre l'audit de sécurité complet
type AuditService struct {
	runtime     ports.ContainerRuntime
	vulnScanner ports.VulnerabilityScanner
	compScanner ports.ComplianceScanner
}

// NewAuditService crée un nouveau service d'audit
func NewAuditService(
	runtime ports.ContainerRuntime,
	vulnScanner ports.VulnerabilityScanner,
	compScanner ports.ComplianceScanner,
) *AuditService {
	return &AuditService{
		runtime:     runtime,
		vulnScanner: vulnScanner,
		compScanner: compScanner,
	}
}

// RunAudit exécute un audit complet.
// onProgress est optionnel (peut être nil).
func (s *AuditService) RunAudit(ctx context.Context, onProgress ProgressCallback) (*domain.AuditReport, error) {
	startTime := time.Now()

	// Helper pour notifier le progrès
	notify := func(curr, total int, msg string) {
		if onProgress != nil {
			onProgress(curr, total, msg)
		}
	}

	notify(0, 0, "Connexion au runtime...")
	if err := s.runtime.Ping(ctx); err != nil {
		return nil, fmt.Errorf("runtime inaccessible: %w", err)
	}

	notify(0, 0, "Listage des conteneurs...")
	containers, err := s.runtime.ListContainers(ctx)
	if err != nil {
		return nil, fmt.Errorf("échec listage conteneurs: %w", err)
	}

	total := len(containers)
	report := &domain.AuditReport{
		Metadata: domain.ReportMetadata{
			GeneratedAt:     time.Now(),
			TotalContainers: total,
		},
		Containers: containers,
		ScanErrors: make([]string, 0),
	}

	// Scan parallèle
	notify(0, total, "Démarrage des scanners...")
	s.scanContainers(ctx, containers, report, notify)

	// Finalisation
	report.Metadata.ScanDuration = time.Since(startTime)
	report.Summary = s.calculateSummary(report)
	report.Summary.RiskScore = report.CalculateRiskScore()

	return report, nil
}

// scanContainers scanne en parallèle avec tolérance aux pannes
func (s *AuditService) scanContainers(ctx context.Context, containers []domain.Container, report *domain.AuditReport, notify ProgressCallback) {
	var (
		wg sync.WaitGroup
		mu sync.Mutex
	)

	// Compteur atomique simulé via mutex pour le progrès
	processedCount := 0

	// Limiteur de concurrence (Worker Pool)
	semaphore := make(chan struct{}, 5) // Max 5 scans simultanés pour ne pas tuer le CPU

	for _, container := range containers {
		wg.Add(1)
		go func(c domain.Container) {
			defer wg.Done()

			// Acquisition du sémaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// 1. Inspection Runtime
			details, err := s.runtime.InspectContainer(ctx, c.ID)

			mu.Lock()
			if err != nil {
				report.ScanErrors = append(report.ScanErrors, fmt.Sprintf("Conteneur %s (%s): %v", c.Name, c.ID, err))
				processedCount++
				notify(processedCount, len(containers), fmt.Sprintf("Erreur sur %s", c.Name))
				mu.Unlock()
				return
			}
			mu.Unlock()

			// 2. Scan Vulnérabilités (Trivy)
			if s.vulnScanner != nil && s.vulnScanner.IsAvailable() {
				// On scanne l'IMAGE du conteneur
				vulns, err := s.vulnScanner.ScanImage(ctx, c.Image)
				mu.Lock()
				if err != nil {
					// On log l'erreur mais on ne fail pas tout le rapport
					report.ScanErrors = append(report.ScanErrors, fmt.Sprintf("Scan Trivy %s: %v", c.Image, err))
				} else {
					report.Vulnerabilities = append(report.Vulnerabilities, vulns...)
				}
				mu.Unlock()
			}

			// 3. Compliance Scan
			if s.compScanner != nil {
				misconfigs, err := s.compScanner.CheckCompliance(ctx, *details)
				mu.Lock()
				if err != nil {
					report.ScanErrors = append(report.ScanErrors, fmt.Sprintf("Compliance check %s: %v", c.Name, err))
				} else {
					report.Misconfigurations = append(report.Misconfigurations, misconfigs...)
				}
				mu.Unlock()
			}

			// Update Progress
			mu.Lock()
			processedCount++
			notify(processedCount, len(containers), fmt.Sprintf("Analysé : %s", c.Name))
			mu.Unlock()

		}(container)
	}

	wg.Wait()
}

func (s *AuditService) calculateSummary(report *domain.AuditReport) domain.AuditSummary {
	summary := domain.AuditSummary{
		TotalVulnerabilities:   len(report.Vulnerabilities),
		TotalMisconfigurations: len(report.Misconfigurations),
	}

	for _, vuln := range report.Vulnerabilities {
		if vuln.IsCritical() {
			if vuln.Severity == domain.SeverityCritical {
				summary.CriticalVulnerabilities++
			} else {
				summary.HighVulnerabilities++
			}
		}
	}

	// Détection naïve des privilèges (à améliorer avec les details réels)
	// Dans une V2, on mapperait les ContainerDetails dans le rapport pour plus de précision
	for _, m := range report.Misconfigurations {
		if m.ID == "PANOPTIC-001" { // ID réservé pour Privileged Container
			summary.PrivilegedContainers++
		}
	}

	return summary
}
