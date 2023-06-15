package repositories

import (
	"context"
	"crypto-satangpro/db"
	"crypto-satangpro/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateTransactionRepo(taskModel models.TransactionModel) (response string, err error) {

	mongoDB, err := db.GetMongoDB()
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mongoDB.Collection("transactions")

	result, err := collection.InsertOne(ctx, taskModel)
	if err != nil {
		return
	}

	response = result.InsertedID.(primitive.ObjectID).Hex()
	return
}