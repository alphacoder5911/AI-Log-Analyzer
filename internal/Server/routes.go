package server

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine,hub *Hub,IncidentHandler *IncidentHandler){

	router.GET("/ws",func(c *gin.Context){
		ServerWs(hub,c)

	})


	router.GET("/incident/recent",IncidentHandler.GetRecentIncidents)
}