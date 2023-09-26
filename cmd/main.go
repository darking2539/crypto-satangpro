package main

import (
	"crypto-satangpro/cmd/worker"
	"crypto-satangpro/db"
	"crypto-satangpro/models"
	"crypto-satangpro/utils"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
	"go.uber.org/ratelimit"
)

func init() {
	godotenv.Load()
	err := db.InitMongoDB()
	if err != nil {
		log.Panicln(err.Error())
		return
	}
}

func main() {

	startTime := time.Now()
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	//init work Group
	workerCount := 5
	batchSize := 10
	workerPool := worker.NewWorkerPool(workerCount, batchSize)
	log.Printf("worker: %d, batchSize: %d", workerCount, batchSize)

	go func() {
		super := <-sigterm
		
		workerPool.Stop()
		endTime := time.Now()
		elapsedTime := endTime.Sub(startTime)
		totalInserted := 0
	
		for _, worker := range workerPool.Workers {
			totalInserted += worker.Inserted
		}
	
		log.Println("TimeTaken:", elapsedTime)
		log.Println("Inserted Record:", totalInserted)
		log.Println(super)
		log.Println("Wait 5s before shut down!!!")
		time.Sleep(time.Second*5)
		os.Exit(0)
	}()

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

	//start feed data to workerPool
	startBlock := uint64(18000500)
	currentBlockInt, err := utils.HexToBigInt(currentBlock)
	log.Printf("currentBlock: %d", currentBlockInt.Uint64())

	for i := startBlock;  i <= currentBlockInt.Uint64(); i++ {
		var totalTransaction string
		var wg sync.WaitGroup
		blockHexString := utils.IntToHex(i)
		err := client.Call(&totalTransaction, "eth_getBlockTransactionCountByNumber", blockHexString)
		if err != nil {
			//log.Println("block: %d, error: %s", i, err.Error())
			continue
		}

		totalTransactionInt, err := utils.HexToBigInt(totalTransaction)

		//log show
		log.Printf("block: %d, totalTransaction: %d", i, totalTransactionInt.Uint64())
		rl := ratelimit.New(100000)
		for j := uint64(0); j < totalTransactionInt.Uint64(); j++ {
			rl.Take()
			wg.Add(1)
			go func(index uint64) {
				defer wg.Done()
				var tx models.TransactionResponse
				err := client.Call(&tx, "eth_getTransactionByBlockNumberAndIndex", blockHexString, utils.IntToHex(index))
	
				if err != nil {
					//log.Printf("block: %d, transaction: %d error: %s", i, index, err.Error())
					return
				}

				bigIntBlockNo, err := utils.HexToBigInt(tx.BlockNo)
				if err != nil {
					log.Printf("HexToBigInt error: %s", err.Error())
					return
				}
			
				bigIntIndex, err := utils.HexToBigInt(tx.TransactionIndex)
				if err != nil {
					log.Printf("HexToBigInt error: %s", err.Error())
					return
				}

				record := models.TransactionModel{
					BlockNo: bigIntBlockNo.Uint64(),
					TransactionIndex: bigIntIndex.Uint64(),
					Hash: tx.Hash,
					From: tx.From,
					To: tx.To,
					Value: tx.Value,
					Gas: tx.Gas,
					GasPrice: tx.GasPrice,
				}

				workerPool.Submit(record)
			}(j)
		}

		//wg.Wait() //wait untils finished all
	}


	workerPool.Stop()
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)
	totalInserted := 0

	for _, worker := range workerPool.Workers {
		totalInserted += worker.Inserted
	}

	log.Println("TimeTaken:", elapsedTime)
	log.Println("Inserted Record:", totalInserted)
	

}