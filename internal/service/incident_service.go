package service

import (
	"errors"
	repository "og_Analyzer/internal/Repository"
	"og_Analyzer/internal/models"
)


type InciServiceInt interface{
	GetRecentIncidents(limit int) (*[]models.Incidents,error)
}

type INCI_SERVICE struct {
	IncidentRepo repository.Incii
}


func NewIncidentService(Repo repository.Incii) *INCI_SERVICE{
	return &INCI_SERVICE{
		IncidentRepo: Repo,
	}
}


func (IS *INCI_SERVICE) GetRecentIncidents(limit int) (*[]models.Incidents,error){

	if limit>20{
		return nil,errors.New("Limit must not cross 20")
	}

	res,err:=IS.IncidentRepo.GetRecentIncidents(limit)
	if err!=nil{
		return nil,err

	}

	return res,nil
}