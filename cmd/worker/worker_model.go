package worker

import (
	"crypto-satangpro/models"
	"sync"
)

type WorkerPool struct {
	Workers     []*Worker
	queue       chan *models.TransactionModel
	stop        chan bool
	wg          sync.WaitGroup
	workersCount int
	batchSize   int
}

type Worker struct {
	queue     chan *models.TransactionModel
	stop      chan bool
	wg        *sync.WaitGroup
	records   []*models.TransactionModel
	batchSize int
	Inserted  int
}