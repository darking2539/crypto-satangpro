package models

import "time"

type UserModel struct {
	Address     string    `bson:"address"`
	CreatedBy   string    `bson:"createdBy"`
	CreatedDate time.Time `bson:"createdDate"`
}
