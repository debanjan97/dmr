package worker

import (
	"log"
	"net/rpc"
	"time"

	"github.com/google/uuid"
)

type Worker struct {
	id   string
	l    log.Logger
	conn *rpc.Client
}

func NewWorker() *Worker {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	logger := log.Default()
	logger.SetPrefix("worker")

	return &Worker{
		id:   uuid.New().String(),
		l:    *log.New(logger.Writer(), logger.Prefix(), logger.Flags()),
		conn: client,
	}
}

func (w *Worker) Register() {
	var isRegistered bool
	for !isRegistered {
		time.Sleep(1 * time.Second)
		w.conn.Call("Coordinator.RegisterWorker", w.id, &isRegistered)
	}
}

func (w *Worker) WaitForInstructions() {
	for {
		w.l.Println("waiting for instructions")
		time.Sleep(15 * time.Second)
		w.conn.Call("Coordinator.AskForInstructions", w.id, map[string]interface{}{})
	}
}
