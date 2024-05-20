package main

import (
	configs "frderoubaix.me/cron-as-a-service/config"
	"frderoubaix.me/cron-as-a-service/config/database"
	"frderoubaix.me/cron-as-a-service/cron"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
}

func main() {
	// Charger les variables d'environnement depuis le fichier .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	database.InitDb()
	configs.InitSentry()
	cron.InitCron()
	configs.InitRouter()
}
