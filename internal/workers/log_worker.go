package workers

import (
	"context"

	"fmt"

	repository "og_Analyzer/internal/Repository"
	server "og_Analyzer/internal/Server"

	"og_Analyzer/internal/service"

	
	"github.com/rs/zerolog/log"
)

type LogWorker struct {
	repo *service.LogRepo 
	AnalyZer *service.AnalyZer
	Hub 	*server.Hub
	IncidentRepo repository.Incii
}

func NewLogWorker(r *service.LogRepo,a *service.AnalyZer,hub *server.Hub , inci repository.Incii) *LogWorker{
	return &LogWorker{
		repo: r,
		AnalyZer: a,
		Hub: hub,	
		IncidentRepo: inci,
	}
}

func( L *LogWorker) StartWorker(ctx context.Context) {
	for {
		logEntry,err:=L.repo.PopLog(ctx)
		if err!=nil{
			fmt.Println("Error popping log",err)
			continue
		}
	fmt.Printf("[Worker] Processed :%s - %s \n",logEntry.Level,logEntry.Message)
		if logEntry.Level=="ERROR"{
			str,err:=L.repo.GetAnalysis(ctx,logEntry.Message)
			if err!=nil{
				fmt.Println("Failed to analyze",err)
			}	
			
			
			fmt.Printf("%+v\n",str)



			result := map[string]interface{}{
        "log":      logEntry,
        "analysis": str,
    }

	
	vall, err := L.IncidentRepo.SaveIncident(ctx, logEntry, str)
	if err != nil {
		log.Error().Err(err).Msg("Failed to save incident")
	}
	log.Info().Interface("Incident",vall).Msg("Pushed incident to db")//recommended way to set logs 
	

    // PUSH TO THE HUB
    L.Hub.Broadcast <- result
		}


		
	}
}
