package rbac

import (
	"context"
	"fmt"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateRoleBinding(clientset *kubernetes.Clientset, namespace, username, role string) error {
	roleBindingName := fmt.Sprintf("%s-%s-binding", username, role)
	roleBinding, err := clientset.RbacV1().RoleBindings(namespace).Get(context.TODO(), roleBindingName, metav1.GetOptions{})
	if err == nil {
		// RoleBinding já existe, vamos atualizá-lo
		roleBinding.RoleRef.Name = role
		_, err = clientset.RbacV1().RoleBindings(namespace).Update(context.TODO(), roleBinding, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("erro ao atualizar RoleBinding: %v", err)
		}
		fmt.Printf("RoleBinding %s atualizado no namespace %s\n", roleBinding.Name, namespace)
		return nil
	}
	if !errors.IsNotFound(err) {
		return fmt.Errorf("erro ao verificar RoleBinding: %v", err)
	}

	roleBinding = &rbacv1.RoleBinding{
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
			Kind:     "ClusterRole",
			Name:     role,
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
	clusterRoleBindingName := fmt.Sprintf("%s-%s-binding", username, role)
	clusterRoleBinding, err := clientset.RbacV1().ClusterRoleBindings().Get(context.TODO(), clusterRoleBindingName, metav1.GetOptions{})
	if err == nil {
		// ClusterRoleBinding já existe, vamos atualizá-lo
		clusterRoleBinding.RoleRef.Name = role
		_, err = clientset.RbacV1().ClusterRoleBindings().Update(context.TODO(), clusterRoleBinding, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("erro ao atualizar ClusterRoleBinding: %v", err)
		}
		fmt.Printf("ClusterRoleBinding %s atualizado\n", clusterRoleBinding.Name)
		return nil
	}
	if !errors.IsNotFound(err) {
		return fmt.Errorf("erro ao verificar ClusterRoleBinding: %v", err)
	}

	clusterRoleBinding = &rbacv1.ClusterRoleBinding{
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
	roleBindingName := fmt.Sprintf("%s-%s-binding", username, role)
	roleBinding, err := clientset.RbacV1().RoleBindings(namespace).Get(context.TODO(), roleBindingName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return CreateRoleBinding(clientset, namespace, username, role)
	}
	if err != nil {
		return fmt.Errorf("erro ao obter RoleBinding: %v", err)
	}

	roleBinding.RoleRef.Name = role
	_, err = clientset.RbacV1().RoleBindings(namespace).Update(context.TODO(), roleBinding, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("erro ao atualizar RoleBinding: %v", err)
	}

	fmt.Printf("RoleBinding %s atualizado no namespace %s\n", roleBinding.Name, namespace)
	return nil
}

func UpdateClusterRoleBinding(clientset *kubernetes.Clientset, username, role string) error {
	clusterRoleBindingName := fmt.Sprintf("%s-%s-binding", username, role)
	clusterRoleBinding, err := clientset.RbacV1().ClusterRoleBindings().Get(context.TODO(), clusterRoleBindingName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return CreateClusterRoleBinding(clientset, username, role)
	}
	if err != nil {
		return fmt.Errorf("erro ao obter ClusterRoleBinding: %v", err)
	}

	clusterRoleBinding.RoleRef.Name = role
	_, err = clientset.RbacV1().ClusterRoleBindings().Update(context.TODO(), clusterRoleBinding, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("erro ao atualizar ClusterRoleBinding: %v", err)
	}

	fmt.Printf("ClusterRoleBinding %s atualizado\n", clusterRoleBinding.Name)
	return nil
}
