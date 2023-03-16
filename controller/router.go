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
		GET("api/k8s/pods", Pod.GetPod)

}
