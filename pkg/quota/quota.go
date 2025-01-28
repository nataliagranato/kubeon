package quota

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateResourceQuota(clientset *kubernetes.Clientset, namespace, name string, limits corev1.ResourceList) error {
	// Verifica se já existe uma ResourceQuota no namespace
	existingQuotas, err := clientset.CoreV1().ResourceQuotas(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("erro ao listar ResourceQuotas: %v", err)
	}

	if len(existingQuotas.Items) > 0 {
		return fmt.Errorf("já existe uma ResourceQuota no namespace %s. Use a opção --action=update para atualizar a quota existente.", namespace)
	}

	resourceQuota := &corev1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: corev1.ResourceQuotaSpec{
			Hard: limits,
		},
	}

	_, err = clientset.CoreV1().ResourceQuotas(namespace).Create(context.TODO(), resourceQuota, metav1.CreateOptions{})
	if err != nil {
		if errors.IsAlreadyExists(err) {
			return fmt.Errorf("resourcequotas \"%s\" já existe. Use a opção --action=update para atualizar a quota existente.", name)
		}
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

func DeleteResourceQuota(clientset *kubernetes.Clientset, namespace, name string) error {
	err := clientset.CoreV1().ResourceQuotas(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("erro ao deletar ResourceQuota: %v", err)
	}

	fmt.Printf("ResourceQuota %s deletada no namespace %s\n", name, namespace)
	return nil
}
