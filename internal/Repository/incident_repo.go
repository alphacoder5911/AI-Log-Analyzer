package repository

import (
	"context"
	"og_Analyzer/internal/models"

	"gorm.io/gorm"
)

type Incii interface {
	SaveIncident(ctx context.Context, log models.LogEntry, analysis models.AIAnalysis) (*models.Incidents, error)
}

type IncidentRepo struct {
	db *gorm.DB
}

func NewIncidentRepo(db *gorm.DB) *IncidentRepo{
	return &IncidentRepo{
		db: db,
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

	if err:=I.db.Create(&Inci).Error; err!=nil{
		return nil,err
	}

		return &Inci,nil
 
}