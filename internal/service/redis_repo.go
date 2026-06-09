package service


import (
	"context"
	"time"

	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"og_Analyzer/internal/models"
	

	//"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)
// THis does the work of redis queuing 
type LogRepo struct {
	client *redis.Client
	Analyzer AnalyZer
}


func NewLogREpo(r *redis.Client,a AnalyZer ) *LogRepo{
	return &LogRepo{
		client: r,
		Analyzer: a,
	}
}

func (L *LogRepo) PushLog(c context.Context,entry models.LogEntry) error{

	data,_:=json.Marshal(entry)

	if entry.Level=="ERROR"{
		return L.client.LPush(c,"log_critical",data).Err()
	}

	return L.client.LPush(c,"log_queue",data).Err()
}

func (L *LogRepo) PopLog(c context.Context) (models.LogEntry,error){

	res,err:=L.client.BRPop(c,0,"log_queue","log_critical").Result()
	if err!=nil{
		return models.LogEntry{},err
	}

	var entry models.LogEntry
	err=json.Unmarshal([]byte(res[1]),&entry)
	return entry,err
}

func (L *LogRepo) GetAnalysis(ctx context.Context,message string)(models.AIAnalysis,error){
	encrypt:=sha256.Sum256([]byte(message))
	hex:=hex.EncodeToString(encrypt[:])
	cacheKey := "ai_cache:" + hex
	result,err:=L.client.Get(ctx,cacheKey).Result()
	
	if (err== redis.Nil){
		fmt.Println("Redis cache is empty,going for Ai service")



		res,err:=L.Analyzer.AnalyzeError(ctx,message)
		if err!=nil{
			return models.AIAnalysis{},err
		}

		cacheData,err:=json.Marshal(res)

			err =L.client.Set(ctx,cacheKey,cacheData,24*time.Hour).Err()
			if err!=nil{
				return models.AIAnalysis{},err
			}
		

		return res,nil 
	}else if (err!=nil){
		return models.AIAnalysis{},err
	}else{
		var Analyse models.AIAnalysis
	err=json.Unmarshal([]byte(result),&Analyse)
		return Analyse ,nil
	}
}