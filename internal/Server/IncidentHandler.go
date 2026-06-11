package server

import (
	"net/http"
	"og_Analyzer/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type IncidentHandler struct {
	InciService service.InciServiceInt
}

func NewIncidentHandler(InService service.InciServiceInt) *IncidentHandler{
	return &IncidentHandler{
		InciService: InService,
	}
}

func (h *IncidentHandler) GetRecentIncidents(c *gin.Context){
	limitStr:=c.DefaultQuery("limit","10")
	limit,err:=strconv.Atoi(limitStr)
	if err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"error":err.Error(),
		})
		log.Error().Interface("Error :",err).Msg("invalid limit")
		return
	}

	incidents,err:=h.InciService.GetRecentIncidents(limit)
	if err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"error":err.Error(),
		})
		log.Error().Interface("Error :",err).Msg("Failed at service ")
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"count":len(*incidents),
		"Data":incidents,
	})
}