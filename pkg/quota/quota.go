package quota

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateResourceQuota(clientset *kubernetes.Clientset, namespace, name string, limits corev1.ResourceList) error {
	resourceQuota := &corev1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: corev1.ResourceQuotaSpec{
			Hard: limits,
		},
	}

	_, err := clientset.CoreV1().ResourceQuotas(namespace).Create(context.TODO(), resourceQuota, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("erro ao criar ResourceQuota: %v", err)
	}

	fmt.Printf("ResourceQuota %s criada no namespace %s\n", name, namespace)
	return nil
}

func UpdateResourceQuota(clientset *kubernetes.Clientset, namespace, name string, limits corev1.ResourceList) error {
	resourceQuota, err := clientset.CoreV1().ResourceQuotas(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("erro ao obter ResourceQuota: %v", err)
	}

	resourceQuota.Spec.Hard = limits
	_, err = clientset.CoreV1().ResourceQuotas(namespace).Update(context.TODO(), resourceQuota, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("erro ao atualizar ResourceQuota: %v", err)
	}

	fmt.Printf("ResourceQuota %s atualizada no namespace %s\n", name, namespace)
	return nil
}
