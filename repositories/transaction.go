package repositories

import (
	"context"
	"crypto-satangpro/db"
	"crypto-satangpro/models"
	"math"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func GetTransactionListRepo(page int64, perPage int64, address string) (response []models.TransactionModel, pagination models.Pagination, err error) {

	mongoDB, err := db.GetMongoDB()
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := mongoDB.Collection("transactions")

	filter := bson.M{}
	filter["$or"] = bson.A{
		bson.M{"from": address},
		bson.M{"to": address},
	}

	//pagination
	limit := perPage
	skip := (page - 1) * limit


	sortIndex := bson.D{
		{Key: "blockNumber", Value: -1},
		{Key: "transactionIndex", Value: 1},
	}

	findOptions := options.FindOptions{
		Limit: &limit,
		Skip:  &skip,
		Sort:  &sortIndex,
	}

	totalResults, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return
	}

	cursor, err := collection.Find(ctx, filter, &findOptions)
	if err != nil {
		return
	}

	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &response); err != nil {
		return
	}

	pagination = models.Pagination{
		Page: page,
		PerPage: perPage,
		TotalResults: totalResults,
		TotalPages: int64(math.Ceil(float64(totalResults)/float64(perPage))),
	}

	return
}