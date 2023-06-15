package scheduler

import (
	"crypto-satangpro/models"
	"crypto-satangpro/rabbitmq"
	"crypto-satangpro/utils"
	"os"
	"sync"

	"log"

	"github.com/ethereum/go-ethereum/rpc"
)

var currentBlockInSystem uint64 = 0

func FetchBlockJob() {

	rpcUrl := os.Getenv("ETH_JSON_RPC_API_URL")

	client, err := rpc.Dial(rpcUrl)
	if err != nil {
		log.Panicln(err.Error())
		return
	}

	//get lastest block No From API
	var currentBlock string
	err = client.Call(&currentBlock, "eth_blockNumber")
	if err != nil {
		log.Panicln(err.Error())
		return
	}

	blockNo := currentBlock
	var totalTransaction string
	err = client.Call(&totalTransaction, "eth_getBlockTransactionCountByNumber", blockNo)
	if err != nil {
		log.Panicln(err.Error())
		return
	}

	currentBlockInt, err := utils.HexToBigInt(currentBlock)
	if err != nil {
		log.Panicln(err.Error())
		return
	}

	totalTransactionInt, err := utils.HexToBigInt(totalTransaction)
	if err != nil {
		log.Panicln(err.Error())
		return
	}

	if currentBlockInt.Uint64() > currentBlockInSystem {
		currentBlockInSystem = currentBlockInt.Uint64()
	} else {
		//this block fetch already
		return
	}

	SendAllDataToRabbit(totalTransactionInt.Uint64(), blockNo, client)
}

func SendAllDataToRabbit(totalTransaction uint64, blockNoHexString string, client *rpc.Client) {

	var wg sync.WaitGroup

	conn, ch, q, err := rabbitmq.InitToSendData()
	if err != nil {
		log.Panicln(err.Error())
		return
	}
	defer conn.Close()
	defer ch.Close()

	for i := uint64(0); i < totalTransaction; i++ {
		wg.Add(1)
		go func(index uint64) {
			defer wg.Done()
			var tx models.TransactionResponse
			err := client.Call(&tx, "eth_getTransactionByBlockNumberAndIndex", blockNoHexString, utils.IntToHex(index))

			if err != nil {
				log.Panicln(err.Error())
			}

			rabbitmq.SendingData(tx, conn, ch, *q)
		}(i)
	}

	wg.Wait()

	bigInt, _ := utils.HexToBigInt(blockNoHexString)
	log.Println("Sucessful Block:", bigInt.Uint64(), "total Transaction:", totalTransaction)
}
