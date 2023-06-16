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
	
	InitlizeAddressToDB()
	rabbitmq.InitSubcribeCode(RecordDataService)

}

func RecordDataService(bodyResp []byte) {

	transactionData := models.TransactionResponse{}
	err := json.Unmarshal(bodyResp, &transactionData)
	if err != nil {
		log.Println(err.Error())
		return
	}

	//check address to monitoring
	addressMonitor, repoErr := repositories.CheckUserExistsRepo(transactionData.From, transactionData.To)
	if repoErr != nil {
		fmt.Println("abosszzzz222")
		log.Panicln(err.Error())
		return
	}

	if addressMonitor == "" {
		//address not exist in array
		return
	}

	bigIntBlockNo, err := utils.HexToBigInt(transactionData.BlockNo)
	if err != nil {
		log.Panicln(err.Error())
		return
	}

	bigIntIndex, err := utils.HexToBigInt(transactionData.TransactionIndex)
	if err != nil {
		log.Panicln(err.Error())
		return
	}

	//prepare Data to save db
	dao := models.TransactionModel{
		BlockNo: bigIntBlockNo.Uint64(),
		TransactionIndex: bigIntIndex.Uint64(),
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

	
	payload := fmt.Sprintf("\nAddress: %s\n\nBlockNo: %d\nFrom: %s\nTo: %s\nvalue: %s\nGas: %s", addressMonitor, bigIntBlockNo.Uint64(), transactionData.From, transactionData.To, transactionData.Value, transactionData.Gas)
	utils.LineNotify(payload)

	fmt.Println("record sucessful")	
}

func InitlizeAddressToDB() {
	
	addressString := os.Getenv("ADDRESS_MONITOR")
	addressArray := strings.Split(addressString, ",")

	for _, address := range addressArray {
		
		lowerAddress := strings.ToLower(address)

		dao := models.UserModel{
			Address: lowerAddress,
			CreatedBy: "Initialize",
			CreatedDate: time.Now(),
		}
	
		_, err := repositories.CreateUserRepo(dao)
		if err != nil {
			log.Println("this address registed already", err.Error())
		}
	}
}