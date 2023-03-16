package service

import (
	"github.com/wonderivan/logger"
	"k8s-platform/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var K8s k8s

type k8s struct {
	ClientSet *kubernetes.Clientset
}

func (k *k8s) Init() {
	config, err := clientcmd.BuildConfigFromFlags("", config.Kubeconfig)
	if err != nil {
		logger.Error("创建k8s配置失败，" + err.Error())
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Error("创建k8s ClientSet失败，" + err.Error())
	} else {
		logger.Info("创建k8s ClientSet成功")
	}
	k.ClientSet = clientSet
}
