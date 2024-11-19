package utils

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func BuildKubeConfig(kubeConfigPath string) (*rest.Config, error) {
	if kubeConfigPath != "" {
		cfg, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			return nil, err
		}
		return cfg, nil
	}

	cfg, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
