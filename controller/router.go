package controller

import (
	"github.com/gin-gonic/gin"
)

var Router router

type router struct {
}

func (r *router) InitApiRouter(router *gin.Engine) {
	router.
		GET("/testapi", TestRespon).
		GET("api/k8s/pods", Pod.GetPod).
		GET("api/k8s/pod/detail", Pod.GetPodDetail).
		DELETE("api/k8s/pod/del", Pod.DeletePod).
		PUT("api/k8s/pod/update", Pod.UpdatePod).
		GET("api/k8s/pod/container", Pod.GetPodContainer).
		GET("api/k8s/pod/log", Pod.GetPodLog).
		GET("api/k8s/pod/nspodnum", Pod.GetPodNumPerNp)

}
