package kubeconfig

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func CreateKubeconfig(username string) error {
	// Carregar o kubeconfig existente
	kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.LoadFromFile(kubeconfigPath)
	if err != nil {
		return fmt.Errorf("erro ao carregar kubeconfig existente: %v", err)
	}

	// Obter o cluster atual
	currentContext := config.CurrentContext
	currentCluster := config.Clusters[currentContext]
	if currentCluster == nil {
		return fmt.Errorf("cluster atual não encontrado no kubeconfig")
	}

	// Criar novo config
	newConfig := api.NewConfig()

	// Copiar dados do cluster
	newConfig.Clusters[currentContext] = currentCluster

	// Copiar credenciais do usuário atual
	currentAuthInfo := config.AuthInfos[config.Contexts[currentContext].AuthInfo]
	if currentAuthInfo == nil {
		return fmt.Errorf("credenciais do usuário atual não encontradas")
	}

	// Criar novo authInfo com as mesmas credenciais
	newConfig.AuthInfos[username] = &api.AuthInfo{
		ClientCertificateData: currentAuthInfo.ClientCertificateData,
		ClientKeyData:         currentAuthInfo.ClientKeyData,
		TokenFile:             currentAuthInfo.TokenFile,
		Token:                 currentAuthInfo.Token,
	}

	// Criar novo contexto
	newConfig.Contexts[username] = &api.Context{
		Cluster:  currentContext,
		AuthInfo: username,
	}

	// Definir contexto atual
	newConfig.CurrentContext = username

	// Salvar novo kubeconfig
	newKubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", fmt.Sprintf("config-%s", username))
	err = clientcmd.WriteToFile(*newConfig, newKubeconfigPath)
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
