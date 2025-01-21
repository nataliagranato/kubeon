package cmd

import (
    "fmt"
    "log"

    "github.com/spf13/cobra"
    "kubeon/pkg/kubeconfig"
    "kubeon/utils"
)

var createKubeconfigCmd = &cobra.Command{
    Use:   "create-kubeconfig [username]",
    Short: "Cria um kubeconfig para um novo usuário",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        username := args[0]
        clusterInfo, err := utils.GetClusterInfo()
        if err != nil {
            log.Fatalf("Erro ao obter informações do cluster: %v", err)
        }
        err = kubeconfig.CreateKubeconfig(username, clusterInfo.ClusterName, clusterInfo.Server, clusterInfo.CAData, clusterInfo.ClientCertData, clusterInfo.ClientKeyData)
        if err != nil {
            log.Fatalf("Erro ao criar kubeconfig: %v", err)
        }
        fmt.Printf("Kubeconfig criado para o usuário %s\n", username)
    },
}