package main

import (
	configs "frderoubaix.me/cron-as-a-service/config"
	"frderoubaix.me/cron-as-a-service/config/database"
	"frderoubaix.me/cron-as-a-service/cron"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"log"
	"os"
	"path/filepath"
)

func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
}

func main() {
	// Vérifier l'existence du fichier .env
	envPath := filepath.Join(".", ".env")
	if _, err := os.Stat(envPath); err == nil {
		// Charger les variables d'environnement depuis le fichier .env si le fichier existe
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}
	database.InitDb()
	configs.InitSentry()
	cron.InitCron()
	configs.InitRouter()
}
