package cmd

import (
	"fmt"
	"log"

	"kubeon/pkg/kubeconfig"
	"kubeon/pkg/rbac"
	"kubeon/utils"

	"github.com/spf13/cobra"
)

var role string
var namespace string

var createKubeconfigCmd = &cobra.Command{
	Use:   "create-kubeconfig [username]",
	Short: "Cria um kubeconfig para um novo usuário",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		clientset, err := utils.GetClientset()
		if err != nil {
			log.Fatalf("Erro ao criar clientset: %v", err)
		}

		err = kubeconfig.CreateKubeconfig(username)
		if err != nil {
			log.Fatalf("Erro ao criar kubeconfig: %v", err)
		}

		if role != "" {
			switch role {
			case "admin", "edit", "view":
				if namespace == "" {
					log.Fatalf("Namespace não especificado. Use a flag --namespace para especificar o namespace.")
				}
				err = rbac.CreateRoleBinding(clientset, namespace, username, role)
				if err != nil {
					log.Fatalf("Erro ao criar RoleBinding: %v", err)
				}
				fmt.Printf("Permissões RBAC concedidas ao usuário %s com a role %s no namespace %s\n", username, role, namespace)
			case "cluster-admin":
				err = rbac.CreateClusterRoleBinding(clientset, username, role)
				if err != nil {
					log.Fatalf("Erro ao criar ClusterRoleBinding: %v", err)
				}
				fmt.Printf("Permissões RBAC concedidas ao usuário %s com a role %s no cluster\n", username, role)
			default:
				log.Fatalf("Role inválida: %s. Use 'admin', 'edit', 'view' ou 'cluster-admin'.", role)
			}
		}
	},
}

func init() {
	createKubeconfigCmd.Flags().StringVar(&namespace, "namespace", "", "Namespace para o RoleBinding (obrigatório para roles admin, edit, view)")
	createKubeconfigCmd.Flags().StringVar(&role, "role", "", "Role a ser atribuída ao usuário (admin, edit, view, cluster-admin)")
	rootCmd.AddCommand(createKubeconfigCmd)
}
