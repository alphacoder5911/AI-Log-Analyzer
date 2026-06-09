package server

import (
	

	"github.com/gin-gonic/gin"
)

func StartServer(hub *Hub) *gin.Engine{
	router:=gin.New()

	router.GET("/ws",func(c *gin.Context){
		ServerWs(hub,c)

	})

	return router
}