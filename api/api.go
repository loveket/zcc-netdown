package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zcc_netdown/control"
)

func Api(r *gin.Engine) {
	r.GET("/", control.Getsorcelist)
	r.GET("/downcommonsource", control.DownCommonSource)
	r.POST("/upcommonsource", control.UpCommonSource)
}
