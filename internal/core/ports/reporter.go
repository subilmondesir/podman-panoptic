// internal/core/ports/reporter.go
package ports

import (
	"io"

	"github.com/subilmondesir/podman-panoptic/internal/core/domain"
)

// Reporter définit l'interface pour générer des rapports
type Reporter interface {
	// Generate produit le rapport dans le format cible
	Generate(report *domain.AuditReport, output io.Writer) error
}

// ReportFormat représente les formats de sortie supportés
type ReportFormat string

const (
	FormatJSON ReportFormat = "json"
	FormatHTML ReportFormat = "html"
	FormatText ReportFormat = "text"
)
