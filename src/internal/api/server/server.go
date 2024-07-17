package server

import (
	"github.com/common-nighthawk/go-figure"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	configs "kudejen/src/internal/config"
	"log"
)

type Server struct {
	KubernetesClient *kubernetes.Clientset
}

func NewServer(configPath, kubernetesConfigPath string) *Server {
	BindConfig(configPath)
	initializeServer()
	kubernetesClient := GetKubernetesClient(kubernetesConfigPath)
	return &Server{KubernetesClient: kubernetesClient}
}

func BindConfig(path string) {
	config, err := configs.NewConfig(path)
	if err != nil {
		log.Fatalf("Failed to build config: %v\n", err)
	}
	configs.AppConfig = config
}

var initializeServer = func() {
	myFigure := figure.NewColorFigure("KudeJen", "puffy", "purple", true)
	myFigure.Print()
}

func GetKubernetesClient(kubeConfigPath string) *kubernetes.Clientset {

	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		log.Fatalf("Failed to build kubernetes config: %v\n", err)
		return nil
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create clientset: %v\n", err)
		return nil
	}
	return clientSet
}
