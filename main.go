package main

import (
	"github.com/gin-gonic/gin"
	"k8s-platform/config"
	"k8s-platform/controller"
	"k8s-platform/service"
)

func main() {
	r := gin.Default()
	service.K8s.Init()
	controller.Router.InitApiRouter(r)
	r.Run(config.ListenAddr)
}
