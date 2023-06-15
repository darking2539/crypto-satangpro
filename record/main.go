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

	addressMap := utils.CheckAddressIsExists(transactionData.From, transactionData.To, addressArray)
	
	if addressMap == "" {
		//address not exist in array
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

	
	payload := fmt.Sprintf("\nAddress: %s\n\nBlockNo: %d\nFrom: %s\nTo: %s\nvalue: %s\nGas: %s", addressMap, bigIntBlockNo.Uint64(), transactionData.From, transactionData.To, transactionData.Value, transactionData.Gas)
	utils.LineNotify(payload)

	fmt.Println("record sucessful")
	
}