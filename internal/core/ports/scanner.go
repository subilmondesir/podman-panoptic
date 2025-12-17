// internal/core/ports/scanner.go
package ports

import (
	"context"

	"github.com/subilmondesir/podman-panoptic/internal/core/domain"
)

// VulnerabilityScanner définit l'interface pour scanner les vulnérabilités
type VulnerabilityScanner interface {
	// ScanImage analyse une image de conteneur
	ScanImage(ctx context.Context, imageName string) ([]domain.Vulnerability, error)

	// IsAvailable vérifie si le scanner est prêt à être utilisé
	IsAvailable() bool
}

// ComplianceScanner définit l'interface pour vérifier la conformité
type ComplianceScanner interface {
	// CheckCompliance vérifie les configurations de sécurité
	CheckCompliance(ctx context.Context, container domain.ContainerDetails) ([]domain.Misconfiguration, error)
}
