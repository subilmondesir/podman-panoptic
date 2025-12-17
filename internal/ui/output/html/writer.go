package html

import (
	"embed"
	"html/template"
	"io"
	"strings"

	"github.com/subilmondesir/podman-panoptic/internal/core/domain"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

// Writer génère un rapport HTML sécurisé
type Writer struct {
	tmpl *template.Template
}

// NewWriter crée un nouveau writer HTML
func NewWriter() (*Writer, error) {
	// Fonction helper pour lowercase (utilisée dans le template)
	funcMap := template.FuncMap{
		"lower": strings.ToLower,
	}

	// Parse du template embarqué avec escape automatique XSS
	tmpl, err := template.New("report.tmpl").
		Funcs(funcMap).
		ParseFS(templateFS, "templates/report.tmpl")

	if err != nil {
		return nil, err
	}

	return &Writer{tmpl: tmpl}, nil
}

// Write génère le rapport HTML de manière sécurisée
func (w *Writer) Write(report *domain.AuditReport, output io.Writer) error {
	// Execute sanitize automatiquement toutes les valeurs (anti-XSS)
	return w.tmpl.Execute(output, report)
}
