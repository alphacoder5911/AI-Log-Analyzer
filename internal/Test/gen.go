package main

import (
	"context"
	"fmt"
	"math/rand"
	repository "og_Analyzer/internal/Repository"
	server "og_Analyzer/internal/Server"
	"og_Analyzer/internal/database"
	"og_Analyzer/internal/logger"
	"og_Analyzer/internal/models"
	"og_Analyzer/internal/service"
	"og_Analyzer/internal/workers"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
		log := logger.New()

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading env")
		return
	}

	levels := []string{"ERROR", "INFO"}
	analyzer := service.NewAnalyZer(os.Getenv("GEMINI_API_KEY"))
	
	hub := server.NewHub()
	go hub.Run()

	

	//Initializing database 
	dbUrl:=os.Getenv("DATABASE_URL")
	db,err:=database.Init(dbUrl)
	if err!=nil{
		log.Error().Msg(err.Error())
		return 
		
	}

	mainDb,err:=db.DB()
	if err!=nil{
		log.Fatal().Err(err).Msg("Error loading main db")
		return 
	}

	defer func(){
		if err:=mainDb.Close();err!=nil{
			log.Fatal().Err(err).Msg("Failed to close db")
			
		}
	}()



	IncidentRepo := repository.NewIncidentRepo(db)
	logRepo := service.NewLogREpo(rdb, *analyzer)
	worker := workers.NewLogWorker(logRepo, analyzer, hub, IncidentRepo)

	// 1. Fire up background data processor thread
	go worker.StartWorker(context.Background()) 

	// 2. Wrap log generator loop in a Goroutine so it streams concurrently
	go func() {
		// Small intentional warm up delay to give Gin a split second to claim port :8080
		time.Sleep(2 * time.Second)
		fmt.Println("[Generator] Traffic simulator engine initialized.")

		var v models.LogEntry
		for i := 0; i < 10; i++ {
			time.Sleep(5* time.Second)
			rindex := rand.Intn(len(levels))

			if levels[rindex] == "ERROR" {
				v = models.LogEntry{
					ID:        strconv.Itoa(i),
					Service:   "xyz",
					Level:     levels[rindex],
					Message:   "XYZ service failed due to container issue",
					Timestamp: "00:00",
				}
			} else {
				v = models.LogEntry{
					ID:        strconv.Itoa(i),
					Service:   "xyz",
					Level:     levels[rindex],
					Message:   "XYZ service info",
					Timestamp: "00:00",
				}
			}

			err := logRepo.PushLog(context.Background(), v)
			if err != nil {
				fmt.Println("[Generator] Error pushing log:", err)
				return 
			}
		}
		fmt.Println("[Generator] Run complete. Finished sending 10 test vectors.")
	}()

	// 3. Launch HTTP server infrastructure (Blocks execution and keeps main open)
	serverr := server.StartServer(hub)
	serverr.Run(":8080")
}