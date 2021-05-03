package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd := initRoot()
	rootCmd.AddCommand(initDoc(rootCmd))

	cobra.CheckErr(rootCmd.Execute())

	log.Println("Sniffer is stopped.")
}
