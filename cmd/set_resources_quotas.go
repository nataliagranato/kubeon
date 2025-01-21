package cmd

import (
    "fmt"
    "log"

    "github.com/spf13/cobra"
    corev1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/api/resource"
    "kubeon/pkg/quota"
    "kubeon/utils"
)

var setResourceQuotaCmd = &cobra.Command{
    Use:   "set-resource-quota [namespace]",
    Short: "Define resource quota para um namespace",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        namespace := args[0]
        clientset, err := utils.GetClientset()
        if err != nil {
            log.Fatalf("Erro ao criar clientset: %v", err)
        }
        limits := corev1.ResourceList{
            corev1.ResourceRequestsCPU:    resource.MustParse("1"),
            corev1.ResourceLimitsCPU:      resource.MustParse("2"),
            corev1.ResourceRequestsMemory: resource.MustParse("1Gi"),
            corev1.ResourceLimitsMemory:   resource.MustParse("2Gi"),
        }
        err = quota.CreateResourceQuota(clientset, namespace, "minha-quota", limits)
        if err != nil {
            log.Fatalf("Erro ao criar ResourceQuota: %v", err)
        }
        fmt.Printf("Resource quota definida para o namespace %s\n", namespace)
    },
}