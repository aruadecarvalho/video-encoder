package main

import (
	"log"
	"os"
	"strconv"
	"video-encoder/application/services"
	"video-encoder/framework/database"
	"video-encoder/framework/queue"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

var db database.Database

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	autoMigrateDb, err := strconv.ParseBool(os.Getenv("AUTO_MIGRATE_DB"))
	if err != nil {
		log.Fatalf("Error parsing boolean AUTO_MIGRATE_DB env var")
	}

	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if err != nil {
		log.Fatalf("Error parsing boolean DEBUG env var")
	}

	db.AutoMigrateDb = autoMigrateDb
	db.Debug = debug
	db.DsnTest = os.Getenv("DSN_TEST")
	db.Dsn = os.Getenv("DSN")
	db.DbTypeTest = os.Getenv("DB_TYPE_TEST")
	db.DbType = os.Getenv("DB_TYPE")
	db.Env = os.Getenv("ENV")
}

func main() {
	messageChannel := make(chan amqp.Delivery)
	jobReturnChannel := make(chan services.JobWorkerResult)

	dbConnection, err := db.Connect()
	if err != nil {
		log.Fatalf("Error connecting to DB")
	}

	defer dbConnection.Close()

	rabbitMQ := queue.NewRabbitMQ()
	ch := rabbitMQ.Connect()

	defer ch.Close()

	rabbitMQ.Consume(messageChannel)

	jobManager := services.NewJobManager(dbConnection, rabbitMQ, jobReturnChannel, messageChannel)
	jobManager.Start(ch)
}
