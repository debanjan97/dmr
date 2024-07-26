package main

import (
	"flag"
	"log"
	"time"

	"github.com/debanjan97/distributed_map_reduce/coordinator"
	"github.com/debanjan97/distributed_map_reduce/worker"
)

func main() {
	mode := flag.String("mode", "w", "choose between (c)oordinator and (w)orker")
	flag.Parse()
	log.Default().Printf("mode: %s", *mode)
	switch *mode {
	case "w":
		startWorker()
	case "c":
		startCoordinator()
	default:
		startWorker()
	}
}

func startCoordinator() {
	c := coordinator.NewCoordinator()
	c.Serve()
	c.Start()
	for !c.Done() {
		time.Sleep(1 * time.Second)
	}
}

func startWorker() {
	w := worker.NewWorker()
	w.Register()
	w.Loop()
}
