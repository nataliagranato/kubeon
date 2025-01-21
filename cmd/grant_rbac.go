package cmd

import (
    "fmt"
    "log"

    "github.com/spf13/cobra"
    "kubeon/pkg/rbac"
    "kubeon/utils"
)

var grantRBACCmd = &cobra.Command{
    Use:   "grant-rbac [username] [role]",
    Short: "Concede permissões RBAC a um usuário",
    Args:  cobra.ExactArgs(2),
    Run: func(cmd *cobra.Command, args []string) {
        username := args[0]
        role := args[1]
        clientset, err := utils.GetClientset()
        if err != nil {
            log.Fatalf("Erro ao criar clientset: %v", err)
        }
        err = rbac.CreateRoleBinding(clientset, "meu-namespace", "meu-rolebinding", username, role)
        if err != nil {
            log.Fatalf("Erro ao criar RoleBinding: %v", err)
        }
        fmt.Printf("Permissões RBAC concedidas ao usuário %s com a role %s\n", username, role)
    },
}