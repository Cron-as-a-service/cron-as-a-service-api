package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"os"
	"time"
)

// Client instance
var (
	Client         *mongo.Client
	Tasks          *mongo.Collection
	Database       *mongo.Database
	TasksResult    *mongo.Collection
	TasksLogs      *mongo.Collection
	TasksTreatment *mongo.Collection
)

func InitDb() *mongo.Client {
	mongoUri := os.Getenv("MONGO_URI")
	mongoDatabase := os.Getenv("MONGO_DATABASE")
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoUri).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		zap.L().Error(fmt.Sprintf("Impossible to initialize MongoDB : %v", err))
		panic(err)
	}

	//ping the database
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		zap.L().Error(fmt.Sprintf("Impossible to ping MongoDB : %v", err))
		panic(err)
	}

	Client = client
	Database = Client.Database(mongoDatabase)
	Tasks = Client.Database(mongoDatabase).Collection("task")
	TasksResult = Client.Database(mongoDatabase).Collection("task_result")
	TasksLogs = Client.Database(mongoDatabase).Collection("task_logs")
	TasksTreatment = Client.Database(mongoDatabase).Collection("task_treatment")

	zap.L().Info(fmt.Sprintf("MongoDB crrectly initialize"))
	return client
}

func CloseDbConnection() error {
	return Client.Disconnect(context.Background())
}

// CreateCollection Créer une collection dans MongoDB
func CreateCollection(collectionName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	opts := options.CreateCollection().SetMaxDocuments(1000)
	err := Database.CreateCollection(ctx, collectionName, opts)
	return err
}

// DropCollection Supprimer une collection dans MongoDB
func DropCollection(db *mongo.Database, collectionName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := Database.Collection(collectionName).Drop(ctx)
	return err
}
