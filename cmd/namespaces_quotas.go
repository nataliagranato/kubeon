package cmd

import (
	"context"
	"fmt"
	"log"

	"kubeon/pkg/quota"
	"kubeon/utils"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		if action == "" {
			log.Fatalf("Ação não especificada. Use a flag --action com 'set' ou 'update'.")
		}

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

		quotaName := "quota-" + uuid.New().String()

		switch action {
		case "set":
			err = quota.CreateResourceQuota(clientset, namespace, quotaName, limits)
			if err != nil {
				if errors.IsAlreadyExists(err) {
					log.Fatalf("Erro ao criar ResourceQuota: resourcequotas \"%s\" já existe. Use a opção --action=update para atualizar a quota existente.", quotaName)
				}
				log.Fatalf("Erro ao criar ResourceQuota: %v", err)
			}
			fmt.Printf("Resource quota definida para o namespace %s com o nome %s\n", namespace, quotaName)
		case "update":
			// Verifica se já existe uma ResourceQuota no namespace
			existingQuotas, err := clientset.CoreV1().ResourceQuotas(namespace).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				log.Fatalf("Erro ao listar ResourceQuotas: %v", err)
			}

			if len(existingQuotas.Items) > 0 {
				// Atualiza a primeira ResourceQuota encontrada
				quotaName = existingQuotas.Items[0].Name
			}

			err = quota.UpdateResourceQuota(clientset, namespace, quotaName, limits)
			if err != nil {
				log.Fatalf("Erro ao atualizar ResourceQuota: %v", err)
			}
			fmt.Printf("Resource quota atualizada para o namespace %s com o nome %s\n", namespace, quotaName)
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
	namespaceQuotasCmd.Flags().StringVar(&action, "action", "", "Ação a ser realizada: set ou update (obrigatório)")
}
