package db

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client
var mongoDB *mongo.Database


func InitMongoDB() (err error) {
	
	// Database Config
	databaseURI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(databaseURI)
	
	if client, err = mongo.NewClient(clientOptions); err != nil {
		log.Println(err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	
	err = client.Connect(ctx)
	if err != nil {
		cancel()
		return
	}
	//Cancel context to avoid memory leak
	defer cancel()

	// Ping our db connection
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Println(err.Error())
		return
	}

	databaseName := os.Getenv("DATABASE_NAME")
	mongoDB = client.Database(databaseName)

	log.Println("Initz DB Sucessful")

	return
}

func GetMongoDB() (db *mongo.Database, err error) {

	if client == nil {
		err = InitMongoDB()
		if err != nil {
			log.Println(err.Error())
			return
		}
	}

	if mongoDB == nil {
		err = InitMongoDB()
		if err != nil {
			log.Println(err.Error())
			return
		}
	}

	err = mongoDB.Client().Ping(context.Background(), readpref.Primary())
	if err != nil {
		err = InitMongoDB()
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
	db = mongoDB
	return
}

func Healthz(c *gin.Context) {

	db, err := GetMongoDB()
	if err != nil {
		c.JSON(400, gin.H{"result": "Connect Unsucessful"})
		return
	}

	client := db.Client()

	err = client.Ping(context.Background(), nil)

	if err != nil {
		c.JSON(400, gin.H{"result": "Connect Unsucessful"})
		return
	}

	c.JSON(200, gin.H{"result": "Connect Sucessful"})
}