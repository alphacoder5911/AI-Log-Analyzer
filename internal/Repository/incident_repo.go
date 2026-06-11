package repository

import (
	"context"
	"encoding/json"
	"og_Analyzer/internal/models"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Incii interface {
	SaveIncident(ctx context.Context, log models.LogEntry, analysis models.AIAnalysis) (*models.Incidents, error)
	LPOP(ctx context.Context) (models.Incidents,error)
	GetRecentIncidents(limit int)(*[]models.Incidents,error)
	PersistIncident(ctx context.Context,Incident models.Incidents) (*models.Incidents,error)
}

type IncidentRepo struct {
	db *gorm.DB
	rc *redis.Client 
}

func NewIncidentRepo(db *gorm.DB,c *redis.Client) *IncidentRepo{
	return &IncidentRepo{
		db: db,
		rc: c,
	}
}

func (I *IncidentRepo) SaveIncident(ctx context.Context,log models.LogEntry,Analysis models.AIAnalysis) (*models.Incidents,error){

	//var Incident *models.Incidents

	Inci:=models.Incidents{
		ServiceName: log.Service,
		Level: log.Level,
		Message: log.Message,
		Severity:Analysis.Severity,
		RootCause: Analysis.RootCause,
		Impact: Analysis.Impact,
		SuggestedFix: Analysis.Fix,
		Confidence:Analysis.Confidence,
		
	}

	data,_:=json.Marshal(&Inci)
	if err:=I.rc.LPush(ctx,"Incident_queue",data).Err();err!=nil{
		return &models.Incidents{},err
	}



		return &Inci,nil
 
}

func(I *IncidentRepo) LPOP(ctx context.Context) (models.Incidents,error){

	res,err:=I.rc.BRPop(ctx,5*time.Second,"Incident_queue").Result()
	if err!=nil{
		return models.Incidents{},err
	}

	var logg models.Incidents

	err= json.Unmarshal([]byte(res[1]),&logg)
	return logg,err
}

func(I *IncidentRepo) PersistIncident(ctx context.Context,Incident models.Incidents) (*models.Incidents,error){
	// data,_:=json.Marshal(&Incident)

	if err:= I.db.Create(&Incident).Error;err!=nil{
		return &models.Incidents{},err
	}
	return &Incident,nil
}

func(I *IncidentRepo) GetRecentIncidents(limit int)(*[]models.Incidents,error){
	var incidents []models.Incidents

	err:=I.db.Order("created_at DESC").Limit(limit).Find(&incidents).Error
	if err!=nil{
		return nil,err
	}

	return &incidents,nil
}