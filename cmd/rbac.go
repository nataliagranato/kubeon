package cmd

import (
	"fmt"
	"log"

	"kubeon/pkg/rbac"
	"kubeon/utils"

	"github.com/spf13/cobra"
)

var rbacNamespace string

var rbacCmd = &cobra.Command{
	Use:   "rbac [username] [role]",
	Short: "Atualiza permissões RBAC de um usuário",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		role := args[1]
		clientset, err := utils.GetClientset()
		if err != nil {
			log.Fatalf("Erro ao criar clientset: %v", err)
		}

		switch role {
		case "admin", "edit", "view":
			if rbacNamespace == "" {
				log.Fatalf("Namespace não especificado. Use a flag --namespace para especificar o namespace.")
			}
			err = rbac.CreateRoleBinding(clientset, rbacNamespace, username, role)
			if err != nil {
				log.Fatalf("Erro ao criar RoleBinding: %v", err)
			}
			fmt.Printf("Permissões RBAC atualizadas para o usuário %s com a role %s no namespace %s\n", username, role, rbacNamespace)
		case "cluster-admin":
			err = rbac.CreateClusterRoleBinding(clientset, username, role)
			if err != nil {
				log.Fatalf("Erro ao criar ClusterRoleBinding: %v", err)
			}
			fmt.Printf("Permissões RBAC atualizadas para o usuário %s com a role %s no cluster\n", username, role)
		default:
			log.Fatalf("Role inválida: %s. Use 'admin', 'edit', 'view' ou 'cluster-admin'.", role)
		}
	},
}

var updateRbacCmd = &cobra.Command{
	Use:   "update-rbac [username] [role]",
	Short: "Atualiza permissões RBAC de um usuário existente",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		role := args[1]
		clientset, err := utils.GetClientset()
		if err != nil {
			log.Fatalf("Erro ao criar clientset: %v", err)
		}

		switch role {
		case "admin", "edit", "view":
			if rbacNamespace == "" {
				log.Fatalf("Namespace não especificado. Use a flag --namespace para especificar o namespace.")
			}
			err = rbac.UpdateRoleBinding(clientset, rbacNamespace, username, role)
			if err != nil {
				log.Fatalf("Erro ao atualizar RoleBinding: %v", err)
			}
			fmt.Printf("Permissões RBAC atualizadas para o usuário %s com a role %s no namespace %s\n", username, role, rbacNamespace)
		case "cluster-admin":
			err = rbac.UpdateClusterRoleBinding(clientset, username, role)
			if err != nil {
				log.Fatalf("Erro ao atualizar ClusterRoleBinding: %v", err)
			}
			fmt.Printf("Permissões RBAC atualizadas para o usuário %s com a role %s no cluster\n", username, role)
		default:
			log.Fatalf("Role inválida: %s. Use 'admin', 'edit', 'view' ou 'cluster-admin'.", role)
		}
	},
}

func init() {
	rbacCmd.Flags().StringVar(&rbacNamespace, "namespace", "", "Namespace para o RoleBinding (obrigatório para roles admin, edit, view)")
	updateRbacCmd.Flags().StringVar(&rbacNamespace, "namespace", "", "Namespace para o RoleBinding (obrigatório para roles admin, edit, view)")
}
