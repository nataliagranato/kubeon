package cmd

import (
	"fmt"
	"log"

	"kubeon/pkg/quota"
	"kubeon/utils"

	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

var (
	limitsCPU      string
	limitsMemory   string
	requestsCPU    string
	requestsMemory string
	action         string
)

var namespaceQuotasCmd = &cobra.Command{
	Use:   "namespace-quotas [namespace]",
	Short: "Define ou atualiza resource quota para um namespace",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		namespace := args[0]
		clientset, err := utils.GetClientset()
		if err != nil {
			log.Fatalf("Erro ao criar clientset: %v", err)
		}
		limits := corev1.ResourceList{
			corev1.ResourceRequestsCPU:    resource.MustParse(requestsCPU),
			corev1.ResourceLimitsCPU:      resource.MustParse(limitsCPU),
			corev1.ResourceRequestsMemory: resource.MustParse(requestsMemory),
			corev1.ResourceLimitsMemory:   resource.MustParse(limitsMemory),
		}

		switch action {
		case "set":
			err = quota.CreateResourceQuota(clientset, namespace, "minha-quota", limits)
			if err != nil {
				log.Fatalf("Erro ao criar ResourceQuota: %v", err)
			}
			fmt.Printf("Resource quota definida para o namespace %s\n", namespace)
		case "update":
			err = quota.UpdateResourceQuota(clientset, namespace, "minha-quota", limits)
			if err != nil {
				log.Fatalf("Erro ao atualizar ResourceQuota: %v", err)
			}
			fmt.Printf("Resource quota atualizada para o namespace %s\n", namespace)
		default:
			log.Fatalf("Ação inválida: %s. Use 'set' ou 'update'.", action)
		}
	},
}

func init() {
	namespaceQuotasCmd.Flags().StringVar(&limitsCPU, "limits-cpu", "2", "Limite de CPU")
	namespaceQuotasCmd.Flags().StringVar(&limitsMemory, "limits-memory", "2Gi", "Limite de memória")
	namespaceQuotasCmd.Flags().StringVar(&requestsCPU, "requests-cpu", "1", "Requisição de CPU")
	namespaceQuotasCmd.Flags().StringVar(&requestsMemory, "requests-memory", "1Gi", "Requisição de memória")
	namespaceQuotasCmd.Flags().StringVar(&action, "action", "set", "Ação a ser realizada: set ou update")
}
