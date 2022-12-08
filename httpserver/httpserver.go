package httpserver

import (
	"github.com/gin-gonic/gin"
	"github.com/zcc_netdown/api"
	"github.com/zcc_netdown/utils"
)

var r *gin.Engine

func StartHttpServer() {
	r = gin.Default()
	r.Use(utils.Core)
	r.Static("/static", "static")
	r.LoadHTMLGlob("view/**/*")
	api.Api(r)
	r.Run(":8887")
}
