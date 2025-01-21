package kubeconfig

import (
    "fmt"
    "os"
    "path/filepath"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/tools/clientcmd/api"
)

func CreateKubeconfig(username, clusterName, server string, caData, clientCertData, clientKeyData []byte) error {
    config := api.NewConfig()
    config.Clusters[clusterName] = &api.Cluster{
        Server:                   server,
        CertificateAuthorityData: caData,
    }
    config.AuthInfos[username] = &api.AuthInfo{
        ClientCertificateData: clientCertData,
        ClientKeyData:         clientKeyData,
    }
    config.Contexts[username] = &api.Context{
        Cluster:  clusterName,
        AuthInfo: username,
    }
    config.CurrentContext = username

    kubeconfigPath := filepath.Join(os.Getenv("HOME"), ".kube", fmt.Sprintf("config-%s", username))
    err := clientcmd.WriteToFile(*config, kubeconfigPath)
    if err != nil {
        return fmt.Errorf("erro ao escrever kubeconfig: %v", err)
    }

    fmt.Printf("Kubeconfig criado para o usuário %s em %s\n", username, kubeconfigPath)
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