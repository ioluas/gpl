package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gpl",
	Short: "Practicing The Go Programming Language book",
	Long:  `Implementing examples/exercises as I learn them from the book`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

var execDir string

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	ex, err := os.Executable()
	if err != nil {
		log.Fatalf("error getting executable: %v", err)
	}
	execDir = filepath.Dir(ex)
}
