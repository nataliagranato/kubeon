package rbac

import (
    "context"
    "fmt"
    "k8s.io/client-go/kubernetes"
    rbacv1 "k8s.io/api/rbac/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateRoleBinding(clientset *kubernetes.Clientset, namespace, name, username, role string) error {
    roleBinding := &rbacv1.RoleBinding{
        ObjectMeta: metav1.ObjectMeta{
            Name:      name,
            Namespace: namespace,
        },
        Subjects: []rbacv1.Subject{
            {
                Kind:     rbacv1.UserKind,
                Name:     username,
                APIGroup: rbacv1.GroupName,
            },
        },
        RoleRef: rbacv1.RoleRef{
            Kind:     "Role",
            Name:     role,
            APIGroup: rbacv1.GroupName,
        },
    }

    _, err := clientset.RbacV1().RoleBindings(namespace).Create(context.TODO(), roleBinding, metav1.CreateOptions{})
    if err != nil {
        return fmt.Errorf("erro ao criar RoleBinding: %v", err)
    }

    fmt.Printf("RoleBinding criado para o usuário %s com a role %s no namespace %s\n", username, role, namespace)
    return nil
}