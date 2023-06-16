package repositories

import (
	"context"
	"crypto-satangpro/db"
	"crypto-satangpro/models"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateUserRepo(userModel models.UserModel) (response string, err error) {

	mongoDB, err := db.GetMongoDB()
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := mongoDB.Collection("users")

	filter := bson.M{
		"address": userModel.Address,
	}
	
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return
	}
	
	if count > 0 {
		err = errors.New("this address already register")
		return
	}

	result, err := collection.InsertOne(ctx, userModel)
	if err != nil {
		return
	}

	response = result.InsertedID.(primitive.ObjectID).Hex()
	return
}

func CheckUserExistsRepo(addressFrom string, addressTo string) (address string, err error) {

	mongoDB, err := db.GetMongoDB()
	if err != nil {
		return 
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	collection := mongoDB.Collection("users")
	filter := bson.M{}
	filter["$or"] = bson.A{
		bson.M{"address": addressFrom},
		bson.M{"address": addressTo},
	}
	
	userModel := models.UserModel{}
	err = collection.FindOne(ctx, filter).Decode(&userModel);
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			err = nil
		}else {
			return
		}
	}

	address = userModel.Address

	
	return
}