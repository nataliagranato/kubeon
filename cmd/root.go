package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kubeon",
	Short: "Kubeon CLI",
	Long:  `Kubeon CLI para gerenciar configurações do Kubernetes.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(createKubeconfigCmd)
	rootCmd.AddCommand(deleteKubeconfigCmd)
	rootCmd.AddCommand(rbacCmd)
	rootCmd.AddCommand(updateRbacCmd)
	rootCmd.AddCommand(namespacesQuotasCmd)
	rootCmd.AddCommand(completionCmd)
}
