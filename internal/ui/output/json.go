package output

import (
	"encoding/json"
	"io"

	"github.com/subilmondesir/podman-panoptic/internal/core/domain"
)

// JSONWriter génère un rapport JSON
type JSONWriter struct{}

// NewJSONWriter crée un nouveau writer JSON
func NewJSONWriter() *JSONWriter {
	return &JSONWriter{}
}

// Write génère le rapport JSON
func (w *JSONWriter) Write(report *domain.AuditReport, output io.Writer) error {
	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	return encoder.Encode(report)
}
