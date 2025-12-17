// internal/ui/cli/version.go
package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd représente la commande version
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "📌 Affiche la version de Panoptic",
	Long:  `Affiche les informations de version, commit Git et date de build.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("👁️  PANOPTIC - The All-Seeing Eye\n\n")
		fmt.Printf("Version:    %s\n", version)
		fmt.Printf("Git Commit: %s\n", commit)
		fmt.Printf("Build Date: %s\n", date)
		fmt.Printf("Go Version: %s\n", "go1.22+")
		fmt.Printf("\n")
		fmt.Printf("Repository: https://github.com/subilmondesir/podman-panoptic\n")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
