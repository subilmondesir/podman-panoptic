package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	// 1. Adapters (Infrastructure)
	"github.com/subilmondesir/podman-panoptic/internal/adapters/podman"
	"github.com/subilmondesir/podman-panoptic/internal/adapters/system"
	"github.com/subilmondesir/podman-panoptic/internal/adapters/trivy"

	// 2. Core (Métier & Modèles)
	"github.com/subilmondesir/podman-panoptic/internal/core/domain"
	"github.com/subilmondesir/podman-panoptic/internal/core/services"

	// 3. UI (Présentation)
	"github.com/subilmondesir/podman-panoptic/internal/ui/output"
	"github.com/subilmondesir/podman-panoptic/internal/ui/output/html"
	"github.com/subilmondesir/podman-panoptic/internal/ui/tui"

	// 4. Librairies externes
	tea "github.com/charmbracelet/bubbletea"
)

var (
	outputFormat string
	outputFile   string
	timeout      int
	socketPath   string // Flag pour custom socket
	useTUI       bool   // Flag pour forcer/désactiver TUI
)

// scanCmd représente la commande scan
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "🔍 Lance un audit de sécurité complet",
	Long: `Analyse tous les conteneurs actifs et arrêtés pour détecter :
  • Vulnérabilités CVE (via Trivy)
  • Configurations dangereuses (conteneurs privilégiés, montages sensibles)
  • Secrets exposés dans les variables d'environnement
  • Violations des best practices CIS`,
	RunE: runScan,
}

func init() {
	rootCmd.AddCommand(scanCmd)

	// Définition des flags
	scanCmd.Flags().StringVarP(&outputFormat, "format", "f", "text", "Format de sortie (text, json, html)")
	scanCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Fichier de sortie (défaut: stdout)")
	scanCmd.Flags().IntVarP(&timeout, "timeout", "t", 60, "Timeout global en secondes")
	scanCmd.Flags().StringVarP(&socketPath, "socket", "s", "", "Chemin spécifique du socket Podman")
	scanCmd.Flags().BoolVar(&useTUI, "tui", true, "Utiliser l'interface graphique terminal (TUI)")

	// Binding Viper (pour permettre la config via fichier yaml)
	viper.BindPFlag("scan.format", scanCmd.Flags().Lookup("format"))
	viper.BindPFlag("scan.timeout", scanCmd.Flags().Lookup("timeout"))
	viper.BindPFlag("podman.socket", scanCmd.Flags().Lookup("socket"))
}

// runScan est le point d'entrée de la logique d'audit
func runScan(cmd *cobra.Command, args []string) error {
	// Création du contexte avec Timeout pour éviter les blocages infinis
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	// --- ÉTAPE 1 : Initialisation de l'Infrastructure (Adapters) ---

	// A. Client Podman
	podmanClient, err := podman.NewClient(socketPath)
	if err != nil {
		return fmt.Errorf("initialisation podman: %w\n💡 Vérifiez que le service Podman tourne (systemctl --user start podman.socket)", err)
	}

	// B. Scanner de Conformité (Interne)
	compScanner := system.NewComplianceInspector()

	// C. Scanner de Vulnérabilités (Trivy Wrapper)
	vulnScanner := trivy.NewScanner()

	// Vérification de la disponibilité de Trivy
	trivyAvailable := vulnScanner != nil && vulnScanner.IsAvailable()
	if !trivyAvailable && useTUI {
		// En mode TUI, on affichera le warning visuellement ou via log avant le lancement
		fmt.Println("⚠️  Scanner Trivy non détecté dans le PATH. L'analyse CVE sera désactivée.")
		time.Sleep(1 * time.Second) // Pause courte pour lecture
	}

	// --- ÉTAPE 2 : Instanciation du Cœur (Service) ---

	auditService := services.NewAuditService(
		podmanClient,
		vulnScanner,
		compScanner,
	)

	// --- ÉTAPE 3 : Exécution (Mode Interactif ou Headless) ---

	var report *domain.AuditReport

	// Mode TUI (Interface Graphique Terminal)
	// Conditions : Flag activé, format texte, et sortie standard (pas de fichier)
	if useTUI && outputFormat == "text" && outputFile == "" {

		// Initialisation du modèle Bubble Tea
		initialModel := tui.NewModel(auditService, ctx)
		program := tea.NewProgram(initialModel)

		// Lancement de l'interface
		finalModel, err := program.Run()
		if err != nil {
			return fmt.Errorf("erreur interface TUI: %w", err)
		}

		// Récupération du rapport depuis le modèle final
		// Note : Cela nécessite que tui.Model expose une méthode GetReport() ou que le champ Report soit exporté.
		// Supposons ici que tui.Model possède la méthode GetReport() *domain.AuditReport
		m, ok := finalModel.(tui.Model)
		if !ok {
			return fmt.Errorf("erreur interne: échec du type assertion sur le modèle TUI")
		}

		report = m.GetReport()
		if report == nil {
			// Si l'utilisateur a annulé (Ctrl+C) ou s'il y a eu une erreur fatale dans le TUI
			return fmt.Errorf("audit annulé ou échoué")
		}

	} else {
		// Mode Headless (CI/CD, Logs, Redirection fichier)
		if verbose {
			fmt.Fprintf(os.Stderr, "⏳ Démarrage de l'audit (Mode Headless)...\n")
		}

		// Lancement direct du service sans callback graphique
		report, err = auditService.RunAudit(ctx, nil)
		if err != nil {
			return fmt.Errorf("échec de l'audit: %w", err)
		}
	}

	// --- ÉTAPE 4 : Génération du Rapport Final (Output) ---

	return generateReport(report)
}

// generateReport sélectionne le bon Writer et écrit le résultat
func generateReport(report *domain.AuditReport) error {
	var writer output.Writer
	var err error

	// Sélection du format
	switch outputFormat {
	case "json":
		writer = output.NewJSONWriter()
	case "html":
		writer, err = html.NewWriter()
		if err != nil {
			return fmt.Errorf("initialisation HTML template: %w", err)
		}
	case "text":
		writer = output.NewTextWriter()
	default:
		return fmt.Errorf("format non supporté: %s (utilisez: text, json, html)", outputFormat)
	}

	// Sélection de la destination (Fichier ou Stdout)
	var dest *os.File
	if outputFile != "" {
		// Ouverture/Création du fichier
		f, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("création fichier de sortie: %w", err)
		}
		defer f.Close()
		dest = f

		// Feedback utilisateur
		if outputFormat != "text" || !useTUI {
			fmt.Fprintf(os.Stderr, "💾 Rapport sauvegardé vers : %s\n", outputFile)
		}
	} else {
		dest = os.Stdout
	}

	// Écriture effective
	if err := writer.Write(report, dest); err != nil {
		return fmt.Errorf("écriture du rapport: %w", err)
	}

	return nil
}
