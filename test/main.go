package main

import (
	"crypto-satangpro/scheduler"
	"crypto-satangpro/utils"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()
	rpcUrl := os.Getenv("ETH_JSON_RPC_API_URL")
	blockNo := utils.IntToHex(17065471)

	client, err := rpc.Dial(rpcUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Get the transaction count for the address
	var countHex string
	err = client.Call(&countHex, "eth_getBlockTransactionCountByNumber", blockNo)
	if err != nil {
		log.Fatal(err)
	}

	// Convert hex transaction count to decimal
	countBigInt, err := utils.HexToBigInt(countHex)
	if err != nil {
		log.Fatal(err)
	}

	scheduler.SendAllDataToRabbit(countBigInt.Uint64(),blockNo, client)

}
