package database

import (
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Init(dbURl string) (*gorm.DB, error) {

	var dsn string

	if dbURl != "" {
		log.Info().Msg("URL Retrieved")
		dsn = dbURl
	} else {
		log.Error().Msg("Couldnt retrieve url")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
if err!=nil{
	
	return nil,err
}
	return db, nil
}
