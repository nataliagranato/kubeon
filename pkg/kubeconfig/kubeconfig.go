package kubeconfig

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func CreateKubeconfig(username string, clientCertData, clientKeyData, caCertData []byte) error {
	// Carregar o kubeconfig existente
	kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.LoadFromFile(kubeconfigPath)
	if err != nil {
		return fmt.Errorf("erro ao carregar kubeconfig existente: %v", err)
	}

	// Criar um novo contexto para o usuário
	clusterName := config.CurrentContext
	cluster := config.Clusters[clusterName]
	if cluster == nil {
		return fmt.Errorf("cluster %s não encontrado no kubeconfig existente", clusterName)
	}

	authInfo := &api.AuthInfo{
		ClientCertificateData: clientCertData,
		ClientKeyData:         clientKeyData,
	}
	context := &api.Context{
		Cluster:  clusterName,
		AuthInfo: username,
	}

	// Adicionar o novo contexto ao kubeconfig
	config.AuthInfos[username] = authInfo
	config.Contexts[username] = context
	config.CurrentContext = username

	// Salvar o novo kubeconfig
	newKubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", fmt.Sprintf("config-%s", username))
	err = clientcmd.WriteToFile(*config, newKubeconfigPath)
	if err != nil {
		return fmt.Errorf("erro ao escrever kubeconfig: %v", err)
	}

	fmt.Printf("Kubeconfig criado para o usuário %s em %s\n", username, newKubeconfigPath)
	return nil
}

func DeleteKubeconfig(username string) error {
	kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", fmt.Sprintf("config-%s", username))
	err := os.Remove(kubeconfigPath)
	if err != nil {
		return fmt.Errorf("erro ao remover kubeconfig: %v", err)
	}

	fmt.Printf("Kubeconfig removido para o usuário %s\n", username)
	return nil
}
