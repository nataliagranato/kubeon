package utils

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type ClusterInfo struct {
	ClusterName    string
	Server         string
	CAData         []byte
	ClientCertData []byte
	ClientKeyData  []byte
}

func GetClientset() (*kubernetes.Clientset, error) {
	var kubeconfigPath string
	if home := homedir.HomeDir(); home != "" {
		kubeconfigPath = filepath.Join(home, ".kube", "config")
	} else {
		kubeconfigPath = os.Getenv("KUBECONFIG")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func GetClusterInfo() (*ClusterInfo, error) {
	var kubeconfigPath string
	if home := homedir.HomeDir(); home != "" {
		kubeconfigPath = filepath.Join(home, ".kube", "config")
	} else {
		kubeconfigPath = os.Getenv("KUBECONFIG")
	}

	config, err := clientcmd.LoadFromFile(kubeconfigPath)
	if err != nil {
		return nil, err
	}

	context := config.Contexts[config.CurrentContext]
	cluster := config.Clusters[context.Cluster]
	authInfo := config.AuthInfos[context.AuthInfo]

	return &ClusterInfo{
		ClusterName:    context.Cluster,
		Server:         cluster.Server,
		CAData:         cluster.CertificateAuthorityData,
		ClientCertData: authInfo.ClientCertificateData,
		ClientKeyData:  authInfo.ClientKeyData,
	}, nil
}
