// internal/adapters/system/inspector.go
package system

import (
	"context"
	"fmt"
	"strings"

	"github.com/subilmondesir/podman-panoptic/internal/core/domain"
	"github.com/subilmondesir/podman-panoptic/internal/core/ports"
)

// ComplianceInspector implémente les vérifications de conformité système
type ComplianceInspector struct {
	// Configuration future (seuils, règles personnalisées)
}

// NewComplianceInspector crée un nouvel inspecteur
func NewComplianceInspector() *ComplianceInspector {
	return &ComplianceInspector{}
}

// CheckCompliance vérifie la conformité d'un conteneur
func (i *ComplianceInspector) CheckCompliance(ctx context.Context, container domain.ContainerDetails) ([]domain.Misconfiguration, error) {
	var misconfigs []domain.Misconfiguration

	// Vérification 1: Conteneur privilégié
	if container.IsPrivileged() {
		misconfigs = append(misconfigs, domain.Misconfiguration{
			ID:          "PANOPTIC-001",
			Severity:    domain.SeverityHigh,
			Title:       "Conteneur en mode privilégié",
			Description: fmt.Sprintf("Le conteneur '%s' s'exécute avec --privileged, donnant un accès root complet au système hôte", container.Name),
			Resource:    container.Name,
			Remediation: "Retirer le flag --privileged et utiliser des capabilities Linux spécifiques si nécessaire (--cap-add)",
		})
	}

	// Vérification 2: Montages sensibles
	if container.HasSensitiveMounts() {
		sensitivePaths := i.detectSensitiveMounts(container.Mounts)
		misconfigs = append(misconfigs, domain.Misconfiguration{
			ID:          "PANOPTIC-002",
			Severity:    domain.SeverityHigh,
			Title:       "Montage de répertoires sensibles",
			Description: fmt.Sprintf("Le conteneur '%s' monte des chemins sensibles: %s", container.Name, strings.Join(sensitivePaths, ", ")),
			Resource:    container.Name,
			Remediation: "Limiter les montages aux répertoires strictement nécessaires. Éviter /etc, /var, /root",
		})
	}

	// Vérification 3: Secrets dans les variables d'environnement
	if secrets := i.detectSecretsInEnv(container.EnvironmentVars); len(secrets) > 0 {
		misconfigs = append(misconfigs, domain.Misconfiguration{
			ID:          "PANOPTIC-003",
			Severity:    domain.SeverityCritical,
			Title:       "Secrets potentiels dans les variables d'environnement",
			Description: fmt.Sprintf("Variables suspectes détectées: %s", strings.Join(secrets, ", ")),
			Resource:    container.Name,
			Remediation: "Utiliser des secrets managers (Podman secrets, Vault) au lieu de variables d'environnement",
		})
	}

	// Vérification 4: Network mode host
	if container.NetworkMode == "host" {
		misconfigs = append(misconfigs, domain.Misconfiguration{
			ID:          "PANOPTIC-004",
			Severity:    domain.SeverityMedium,
			Title:       "Mode réseau 'host' utilisé",
			Description: fmt.Sprintf("Le conteneur '%s' partage la stack réseau de l'hôte", container.Name),
			Resource:    container.Name,
			Remediation: "Utiliser un réseau bridge isolé sauf besoin spécifique justifié",
		})
	}

	return misconfigs, nil
}

// detectSensitiveMounts identifie les montages de chemins sensibles
func (i *ComplianceInspector) detectSensitiveMounts(mounts []domain.Mount) []string {
	sensitivePaths := []string{"/etc", "/var", "/root", "/sys", "/proc", "/boot", "/dev"}
	detected := []string{}

	for _, mount := range mounts {
		for _, sensitive := range sensitivePaths {
			if strings.HasPrefix(mount.Source, sensitive) || strings.HasPrefix(mount.Destination, sensitive) {
				detected = append(detected, fmt.Sprintf("%s -> %s", mount.Source, mount.Destination))
				break
			}
		}
	}

	return detected
}

// detectSecretsInEnv détecte les variables d'environnement suspectes
func (i *ComplianceInspector) detectSecretsInEnv(envVars map[string]string) []string {
	suspiciousKeys := []string{
		"PASSWORD", "PASSWD", "PWD",
		"SECRET", "API_KEY", "APIKEY",
		"TOKEN", "AUTH", "PRIVATE_KEY",
		"AWS_SECRET", "DATABASE_PASSWORD",
	}

	detected := []string{}

	for key, value := range envVars {
		keyUpper := strings.ToUpper(key)

		// Vérifier si la clé contient un mot suspect
		for _, suspicious := range suspiciousKeys {
			if strings.Contains(keyUpper, suspicious) && len(value) > 0 {
				detected = append(detected, key)
				break
			}
		}
	}

	return detected
}

// Vérification de l'implémentation de l'interface
var _ ports.ComplianceScanner = (*ComplianceInspector)(nil)
