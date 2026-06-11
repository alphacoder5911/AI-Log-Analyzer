package workers

import (
	"context"
	repository "og_Analyzer/internal/Repository"

	"github.com/redis/go-redis/v9"
	
)
type DB_worker struct {
	repo repository.Incii
	rc *redis.Client
	
	
}

func NewDBWorker(rep repository.Incii,rc *redis.Client) *DB_worker{
	return &DB_worker{
		repo: rep,
		rc: rc,
	}
}

func (D *DB_worker)StartDBWorker(ctx context.Context){
	for{
		res,err:=D.repo.LPOP(ctx)
// 		length, _ := D.rc.LLen(ctx, "Incident_queue").Result()
// log.Info().Int64("queueDepth", length).Msg("Incident queue size")
		if err!=nil{
		// log.Error().Interface("Error",err).Msg("Failed to pop error log ")
		continue

		}

		D.repo.PersistIncident(ctx,res)
		
	}
}