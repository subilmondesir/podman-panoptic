package output

import (
	"io"

	"github.com/subilmondesir/podman-panoptic/internal/core/domain"
)

// Writer définit l'interface pour tous les générateurs de rapports
type Writer interface {
	Write(report *domain.AuditReport, output io.Writer) error
}
