package scheduler

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
)

func RunScheduler() {

	scheduler := gocron.NewScheduler(time.Local)

	_, err := scheduler.Every(10).Seconds().Do(FetchBlockJob)
	if err != nil {
		fmt.Println("Error scheduling task:", err)
		return
	}

	scheduler.StartAsync()

}