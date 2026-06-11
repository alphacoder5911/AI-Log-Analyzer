package server

import (
	

	"github.com/gin-gonic/gin"
)

func StartServer(hub *Hub,IncidentHandler *IncidentHandler) *gin.Engine{
	router:=gin.New()

	RegisterRoutes(router,hub,IncidentHandler)

	return router
}