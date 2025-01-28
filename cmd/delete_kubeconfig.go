package cmd

import (
	"fmt"
	"log"

	"kubeon/pkg/kubeconfig"

	"github.com/spf13/cobra"
)

var deleteKubeconfigCmd = &cobra.Command{
	Use:   "delete-kubeconfig [username]",
	Short: "Remove o kubeconfig de um usuário",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		err := kubeconfig.DeleteKubeconfig(username)
		if err != nil {
			log.Fatalf("Erro ao remover kubeconfig: %v", err)
		}
		fmt.Printf("Kubeconfig removido para o usuário %s\n", username)
	},
}

func init() {
}
