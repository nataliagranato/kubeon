package rbac

import (
	"context"
	"fmt"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func removeExistingRoleBindings(clientset *kubernetes.Clientset, namespace, username string) error {
	roleBindings, err := clientset.RbacV1().RoleBindings(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("erro ao listar RoleBindings: %v", err)
	}

	for _, rb := range roleBindings.Items {
		for _, subject := range rb.Subjects {
			if subject.Kind == "User" && subject.Name == username {
				err = clientset.RbacV1().RoleBindings(namespace).Delete(context.TODO(), rb.Name, metav1.DeleteOptions{})
				if err != nil {
					return fmt.Errorf("erro ao remover RoleBinding existente: %v", err)
				}
			}
		}
	}

	return nil
}

func removeExistingClusterRoleBindings(clientset *kubernetes.Clientset, username string) error {
	clusterRoleBindings, err := clientset.RbacV1().ClusterRoleBindings().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("erro ao listar ClusterRoleBindings: %v", err)
	}

	for _, crb := range clusterRoleBindings.Items {
		for _, subject := range crb.Subjects {
			if subject.Kind == "User" && subject.Name == username {
				err = clientset.RbacV1().ClusterRoleBindings().Delete(context.TODO(), crb.Name, metav1.DeleteOptions{})
				if err != nil {
					return fmt.Errorf("erro ao remover ClusterRoleBinding existente: %v", err)
				}
			}
		}
	}

	return nil
}

func CreateRoleBinding(clientset *kubernetes.Clientset, namespace, username, role string) error {
	err := removeExistingRoleBindings(clientset, namespace, username)
	if err != nil {
		return err
	}

	roleBindingName := fmt.Sprintf("%s-%s-binding", username, role)
	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      roleBindingName,
			Namespace: namespace,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:     "User",
				Name:     username,
				APIGroup: "rbac.authorization.k8s.io",
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole", // Usando ClusterRole padrão do Kubernetes
			Name:     role,          // view, edit, admin ou cluster-admin
			APIGroup: "rbac.authorization.k8s.io",
		},
	}

	_, err = clientset.RbacV1().RoleBindings(namespace).Create(context.TODO(), roleBinding, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("erro ao criar RoleBinding: %v", err)
	}

	fmt.Printf("RoleBinding %s criado no namespace %s\n", roleBinding.Name, namespace)
	return nil
}

func CreateClusterRoleBinding(clientset *kubernetes.Clientset, username, role string) error {
	err := removeExistingClusterRoleBindings(clientset, username)
	if err != nil {
		return err
	}

	clusterRoleBindingName := fmt.Sprintf("%s-%s-binding", username, role)
	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterRoleBindingName,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:     "User",
				Name:     username,
				APIGroup: "rbac.authorization.k8s.io",
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     role,
			APIGroup: "rbac.authorization.k8s.io",
		},
	}

	_, err = clientset.RbacV1().ClusterRoleBindings().Create(context.TODO(), clusterRoleBinding, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("erro ao criar ClusterRoleBinding: %v", err)
	}

	fmt.Printf("ClusterRoleBinding %s criado\n", clusterRoleBinding.Name)
	return nil
}

func UpdateRoleBinding(clientset *kubernetes.Clientset, namespace, username, role string) error {
	err := removeExistingRoleBindings(clientset, namespace, username)
	if err != nil {
		return err
	}

	return CreateRoleBinding(clientset, namespace, username, role)
}

func UpdateClusterRoleBinding(clientset *kubernetes.Clientset, username, role string) error {
	err := removeExistingClusterRoleBindings(clientset, username)
	if err != nil {
		return err
	}

	return CreateClusterRoleBinding(clientset, username, role)
}
