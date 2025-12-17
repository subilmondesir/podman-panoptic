package domain

import "time"

// AuditReport représente le résultat complet d'un audit
type AuditReport struct {
	Metadata          ReportMetadata
	Containers        []Container
	Vulnerabilities   []Vulnerability
	Misconfigurations []Misconfiguration
	ScanErrors        []string // NOUVEAU: Liste des erreurs non-bloquantes (ex: échec d'inspection d'un conteneur)
	Summary           AuditSummary
}

// ReportMetadata contient les métadonnées du rapport
type ReportMetadata struct {
	Version         string
	GeneratedAt     time.Time
	Hostname        string
	TotalContainers int
	ScanDuration    time.Duration
}

// AuditSummary contient les statistiques globales
type AuditSummary struct {
	TotalVulnerabilities    int
	CriticalVulnerabilities int
	HighVulnerabilities     int
	TotalMisconfigurations  int
	PrivilegedContainers    int
	RiskScore               float64 // 0-100
}

// CalculateRiskScore calcule un score de risque global
func (ar *AuditReport) CalculateRiskScore() float64 {
	score := float64(ar.Summary.CriticalVulnerabilities * 10)
	score += float64(ar.Summary.HighVulnerabilities * 5)
	score += float64(ar.Summary.PrivilegedContainers * 3)

	// Pénalité pour les erreurs de scan
	if len(ar.ScanErrors) > 0 {
		score += 5.0
	}

	if score > 100 {
		return 100.0
	}
	return score
}
