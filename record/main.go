package record

import (
	"crypto-satangpro/models"
	"crypto-satangpro/rabbitmq"
	"crypto-satangpro/repositories"
	"crypto-satangpro/utils"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func RunRecord() {
	
	rabbitmq.InitSubcribeCode(RecordDataService)

}

func RecordDataService(bodyResp []byte) {

	transactionData := models.TransactionResponse{}
	err := json.Unmarshal(bodyResp, &transactionData)
	if err != nil {
		log.Panicln(err.Error())
		return
	}

	//get address to monitoring
	
	addressString := os.Getenv("ADDRESS_MONITOR")
	addressArray := strings.Split(addressString, ",")

	if !utils.CheckAddressIsExists(transactionData.From, transactionData.To, addressArray) {
		//if not exist return
		return
	}

	bigIntBlockNo, err := utils.HexToBigInt(transactionData.BlockNo)
	if err != nil {
		log.Panicln(err.Error())
		return
	}

	//prepare Data to save db
	dao := models.TransactionModel{
		BlockNo: bigIntBlockNo.Uint64(),
		Hash: transactionData.Hash,
		From: transactionData.From,
		To: transactionData.To,
		Value: transactionData.Value,
		Gas: transactionData.Gas,
		GasPrice: transactionData.GasPrice,
		CreatedDate: time.Now(),
	}

	_, err = repositories.CreateTransactionRepo(dao)
	if err != nil {
		log.Panicln(err.Error())
		return
	}

	fmt.Println("record sucessful")
	
}