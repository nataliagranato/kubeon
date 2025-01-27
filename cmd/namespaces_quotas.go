package cmd

import (
	"context"
	"fmt"
	"log"

	"kubeon/utils"

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
)

var namespacesQuotasCmd = &cobra.Command{
	Use:   "namespaces-quotas [namespace]",
	Short: "Define quotas de recursos para um namespace",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		namespace := args[0]
		clientset, err := utils.GetClientset()
		if err != nil {
			log.Fatalf("Erro ao criar clientset: %v", err)
		}

		resourceQuota := &corev1.ResourceQuota{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "resource-quota",
				Namespace: namespace,
			},
			Spec: corev1.ResourceQuotaSpec{
				Hard: corev1.ResourceList{
					corev1.ResourceLimitsCPU:      resource.MustParse(limitsCPU),
					corev1.ResourceLimitsMemory:   resource.MustParse(limitsMemory),
					corev1.ResourceRequestsCPU:    resource.MustParse(requestsCPU),
					corev1.ResourceRequestsMemory: resource.MustParse(requestsMemory),
				},
			},
		}

		_, err = clientset.CoreV1().ResourceQuotas(namespace).Create(context.TODO(), resourceQuota, metav1.CreateOptions{})
		if err != nil {
			if errors.IsAlreadyExists(err) {
				_, err = clientset.CoreV1().ResourceQuotas(namespace).Update(context.TODO(), resourceQuota, metav1.UpdateOptions{})
				if err != nil {
					log.Fatalf("Erro ao atualizar ResourceQuota: %v", err)
				}
				fmt.Printf("ResourceQuota atualizada no namespace %s\n", namespace)
			} else {
				log.Fatalf("Erro ao criar ResourceQuota: %v", err)
			}
		} else {
			fmt.Printf("ResourceQuota criada no namespace %s\n", namespace)
		}
	},
}

func init() {
	namespacesQuotasCmd.Flags().StringVar(&limitsCPU, "limits-cpu", "1", "Limite de CPU")
	namespacesQuotasCmd.Flags().StringVar(&limitsMemory, "limits-memory", "1Gi", "Limite de memória")
	namespacesQuotasCmd.Flags().StringVar(&requestsCPU, "requests-cpu", "500m", "Requisição de CPU")
	namespacesQuotasCmd.Flags().StringVar(&requestsMemory, "requests-memory", "512Mi", "Requisição de memória")
	rootCmd.AddCommand(namespacesQuotasCmd)
}
