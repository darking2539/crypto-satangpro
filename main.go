package main

import (
	"crypto-satangpro/db"
	"crypto-satangpro/record"
	"crypto-satangpro/scheduler"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()
	mode := os.Getenv("MODE") //for prod

	if mode == "" {
		mode = os.Args[1] //for dev
	}
	
	if mode == "job" {
		log.Println("run scheduler")
		scheduler.RunScheduler()
		HandleRequests()
	} else if mode == "record" {
		log.Println("run monitoring")
		db.InitMongoDB()
		record.RunRecord()
	} else {
		log.Fatal("mode not found")
	}

}

func HandleRequests() {
	log.Fatal(http.ListenAndServe(":8080", nil))
}