// internal/ui/cli/root.go
package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	verbose bool

	// Version info (injectée au build)
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// rootCmd représente la commande de base appelée sans sous-commandes
var rootCmd = &cobra.Command{
	Use:   "panoptic",
	Short: "👁️  L'œil omniscient de la sécurité des conteneurs",
	Long: `PANOPTIC - The All-Seeing Eye for Container Security

Un outil d'audit de sécurité de nouvelle génération pour environnements 
conteneurisés (Podman/Docker). Analyse en profondeur les vulnérabilités, 
les misconfigurations et les risques de sécurité.

Exemples d'utilisation:
  panoptic scan                    # Audit complet
  panoptic scan --json             # Sortie JSON
  panoptic scan --output report.html  # Rapport HTML
  panoptic version                 # Afficher la version`,

	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute exécute la commande root
func Execute() error {
	return rootCmd.Execute()
}

// SetVersion injecte les informations de version
func SetVersion(v, c, d string) {
	version = v
	commit = c
	date = d
}

func init() {
	cobra.OnInitialize(initConfig)

	// Flags globaux
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Fichier de configuration (default: $HOME/.panoptic.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Mode verbeux (affiche les logs détaillés)")

	// Binding Viper
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

// initConfig lit le fichier de configuration et les variables d'environnement
func initConfig() {
	if cfgFile != "" {
		// Utiliser le fichier de config spécifié
		viper.SetConfigFile(cfgFile)
	} else {
		// Chercher dans le home directory
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".panoptic")
	}

	// Lire les variables d'environnement PANOPTIC_*
	viper.SetEnvPrefix("PANOPTIC")
	viper.AutomaticEnv()

	// Charger le fichier de config (erreur silencieuse si absent)
	if err := viper.ReadInConfig(); err == nil {
		if verbose {
			fmt.Fprintf(os.Stderr, "📝 Configuration chargée: %s\n", viper.ConfigFileUsed())
		}
	}
}
