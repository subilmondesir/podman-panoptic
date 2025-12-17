package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/subilmondesir/podman-panoptic/internal/core/domain"
)

// TextWriter génère un rapport texte formaté
type TextWriter struct{}

// NewTextWriter crée un nouveau writer texte
func NewTextWriter() *TextWriter {
	return &TextWriter{}
}

// Write génère le rapport texte
func (w *TextWriter) Write(report *domain.AuditReport, output io.Writer) error {
	// Header
	fmt.Fprintln(output, strings.Repeat("=", 70))
	fmt.Fprintln(output, "👁️  PANOPTIC - RAPPORT D'AUDIT DE SÉCURITÉ")
	fmt.Fprintln(output, strings.Repeat("=", 70))
	fmt.Fprintln(output, "")

	// Métadonnées
	fmt.Fprintf(output, "Date:         %s\n", report.Metadata.GeneratedAt.Format("2006-01-02 15:04:05"))
	fmt.Fprintf(output, "Version:      %s\n", report.Metadata.Version)
	fmt.Fprintf(output, "Durée:        %s\n", report.Metadata.ScanDuration)
	fmt.Fprintf(output, "Conteneurs:   %d\n", report.Metadata.TotalContainers)
	fmt.Fprintln(output, "")

	// Résumé exécutif
	fmt.Fprintln(output, strings.Repeat("-", 70))
	fmt.Fprintln(output, "📊 RÉSUMÉ EXÉCUTIF")
	fmt.Fprintln(output, strings.Repeat("-", 70))
	fmt.Fprintf(output, "Vulnérabilités:       %d (Critiques: %d, Hautes: %d)\n",
		report.Summary.TotalVulnerabilities,
		report.Summary.CriticalVulnerabilities,
		report.Summary.HighVulnerabilities)
	fmt.Fprintf(output, "Misconfigurations:    %d\n", report.Summary.TotalMisconfigurations)
	fmt.Fprintf(output, "Conteneurs privilégiés: %d\n", report.Summary.PrivilegedContainers)
	fmt.Fprintf(output, "Score de Risque:      %.1f/100\n", report.Summary.RiskScore)

	// Affichage des erreurs de scan éventuelles
	if len(report.ScanErrors) > 0 {
		fmt.Fprintln(output, "")
		fmt.Fprintln(output, "⚠️  ERREURS DE SCAN:")
		for _, err := range report.ScanErrors {
			fmt.Fprintf(output, "  - %s\n", err)
		}
	}
	fmt.Fprintln(output, "")

	// Conteneurs
	fmt.Fprintln(output, strings.Repeat("-", 70))
	fmt.Fprintln(output, "📦 CONTENEURS DÉTECTÉS")
	fmt.Fprintln(output, strings.Repeat("-", 70))

	if len(report.Containers) == 0 {
		fmt.Fprintln(output, "Aucun conteneur détecté")
	} else {
		for _, c := range report.Containers {
			stateIcon := "🟢"
			if c.State != "running" {
				stateIcon = "⚪"
			}
			fmt.Fprintf(output, "%s %-20s | %-15s | %s\n", stateIcon, c.Name, c.State, c.Image)
		}
	}
	fmt.Fprintln(output, "")

	// Misconfigurations
	if len(report.Misconfigurations) > 0 {
		fmt.Fprintln(output, strings.Repeat("-", 70))
		fmt.Fprintln(output, "🛡️  CONFIGURATIONS DE SÉCURITÉ")
		fmt.Fprintln(output, strings.Repeat("-", 70))

		for i, m := range report.Misconfigurations {
			severityIcon := getSeverityIcon(m.Severity)
			fmt.Fprintf(output, "\n[%d] %s %s - %s\n", i+1, severityIcon, m.ID, m.Title)
			fmt.Fprintf(output, "    Ressource:    %s\n", m.Resource)
			fmt.Fprintf(output, "    Sévérité:     %s\n", m.Severity)
			fmt.Fprintf(output, "    Description:  %s\n", m.Description)
			fmt.Fprintf(output, "    Remédiation:  %s\n", m.Remediation)
		}
		fmt.Fprintln(output, "")
	}

	// Footer
	fmt.Fprintln(output, strings.Repeat("=", 70))
	fmt.Fprintln(output, "Fin du rapport")
	fmt.Fprintln(output, strings.Repeat("=", 70))

	return nil
}

// getSeverityIcon retourne l'icône correspondant à la sévérité
func getSeverityIcon(severity domain.Severity) string {
	switch severity {
	case domain.SeverityCritical:
		return "🔴"
	case domain.SeverityHigh:
		return "🟠"
	case domain.SeverityMedium:
		return "🟡"
	case domain.SeverityLow:
		return "🟢"
	default:
		return "ℹ️"
	}
}
