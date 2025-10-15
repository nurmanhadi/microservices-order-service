package worker

import (
	"log"

	"github.com/panjf2000/ants/v2"
)

type WorkerPool struct {
	pool *ants.Pool
}

func NewWorkerPool(size int) *WorkerPool {
	p, err := ants.NewPool(size)
	if err != nil {
		log.Fatalf("failed to create new worker pool: %s", err)
	}
	return &WorkerPool{
		pool: p,
	}
}
