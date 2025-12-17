package tui

import (
	"context"
	"fmt"
	"strings" // Import corrigé (remonté au niveau des imports)

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/subilmondesir/podman-panoptic/internal/core/domain"
	"github.com/subilmondesir/podman-panoptic/internal/core/services"
)

// Styles Lipgloss (Design System)
var (
	titleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#3498db")).Bold(true)
	infoStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#7f8c8d"))
	errStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#e74c3c"))
)

// --- Messages (Events) ---

// scanFinishedMsg indique que l'audit est terminé
type scanFinishedMsg struct {
	Report *domain.AuditReport
	Err    error
}

// progressMsg sert à mettre à jour la barre de progression (pour implémentation future temps réel)
type progressMsg struct {
	Current int
	Total   int
	Message string
}

// --- Modèle Principal (State) ---

type Model struct {
	auditService *services.AuditService
	ctx          context.Context

	// État UI
	loading   bool
	progress  progress.Model
	spinner   spinner.Model
	statusMsg string
	err       error
	report    *domain.AuditReport // Stocke le résultat final

	// Contrôle
	quitting bool
}

// NewModel initialise le modèle Bubble Tea
func NewModel(service *services.AuditService, ctx context.Context) Model {
	// Barre de progression dégradée
	p := progress.New(progress.WithDefaultGradient())
	p.Width = 40

	// Spinner (animation de chargement)
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		auditService: service,
		ctx:          ctx,
		loading:      true,
		progress:     p,
		spinner:      s,
		statusMsg:    "Initialisation de l'audit...",
	}
}

// GetReport permet au CLI de récupérer le résultat une fois le TUI fermé
func (m Model) GetReport() *domain.AuditReport {
	return m.report
}

// Init lance les commandes initiales (Tick spinner + Lancement Audit)
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.runAuditCmd(),
	)
}

// Update gère la boucle d'événements
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Touches Clavier
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}

	// Animation Spinner
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	// Mise à jour Progression (Optionnel V1)
	case progressMsg:
		var cmd tea.Cmd
		if msg.Total > 0 {
			pct := float64(msg.Current) / float64(msg.Total)
			cmd = m.progress.SetPercent(pct)
		}
		m.statusMsg = msg.Message
		return m, cmd

	// Animation Barre de progression
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	// Fin de l'Audit
	case scanFinishedMsg:
		m.loading = false
		if msg.Err != nil {
			m.err = msg.Err
			return m, tea.Quit
		}
		m.report = msg.Report
		return m, tea.Quit // On quitte proprement pour afficher le rapport textuel
	}

	return m, nil
}

// View dessine l'interface dans le terminal
func (m Model) View() string {
	// Cas d'erreur fatale
	if m.err != nil {
		return fmt.Sprintf("\n❌ Erreur fatale: %v\n\n", m.err)
	}

	// Cas de chargement (Scan en cours)
	if m.loading {
		pad := strings.Repeat(" ", padding(0, m.statusMsg))

		return "\n" +
			titleStyle.Render("👁️  PANOPTIC AUDIT") + "\n\n" +
			" " + m.spinner.View() + " " + m.statusMsg + pad + "\n\n" +
			" " + m.progress.View() + "\n\n" +
			infoStyle.Render(" Appuyez sur 'q' pour annuler") + "\n"
	}

	// Une fois fini, on n'affiche rien ici, c'est le CLI TextOutput qui prend le relais
	return ""
}

// --- Commandes ---

// runAuditCmd lance l'audit dans une Goroutine séparée
func (m Model) runAuditCmd() tea.Cmd {
	return func() tea.Msg {
		// Note : Dans cette V1, l'audit est bloquant pour le worker mais non bloquant pour l'UI.
		// Pour une V2 "Real-Time", il faudrait passer un callback au RunAudit qui envoie des progressMsg.

		// Lancement de l'audit
		report, err := m.auditService.RunAudit(m.ctx, nil)

		return scanFinishedMsg{Report: report, Err: err}
	}
}

// --- Helpers ---

func padding(n int, s string) int {
	pad := 40 - len(s)
	if pad < 0 {
		return 1
	}
	return pad
}
