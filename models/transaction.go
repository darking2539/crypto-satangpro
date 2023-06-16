package models

import "time"

type TransactionResponse struct {
	BlockNo          string `json:"blockNumber"`
	TransactionIndex string `json:"transactionIndex"`
	Hash             string `json:"hash"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
}

type TransactionModel struct {
	BlockNo          uint64    `bson:"blockNumber"`
	TransactionIndex uint64    `bson:"transactionIndex"`
	Hash             string    `bson:"hash"`
	From             string    `bson:"from"`
	To               string    `bson:"to"`
	Value            string    `bson:"value"`
	Gas              string    `bson:"gas"`
	GasPrice         string    `bson:"gasPrice"`
	CreatedDate      time.Time `bson:"createdDate"`
}
