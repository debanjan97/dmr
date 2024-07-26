package worker

import (
	"log"
	"net/rpc"
	"time"

	t "github.com/debanjan97/distributed_map_reduce/types"
	"github.com/google/uuid"
)

type Worker struct {
	id   string
	l    log.Logger
	conn *rpc.Client
}

func NewWorker() *Worker {
	logger := log.Default()
	logger.SetPrefix("worker")

	return &Worker{
		id: uuid.New().String(),
		l:  *log.New(logger.Writer(), logger.Prefix(), logger.Flags()),
	}
}

func (w *Worker) Register() {
	for {
		client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
		if err == nil {
			w.conn = client
			break
		}
		w.l.Print("Waiting for controller to come up")
		time.Sleep(1 * time.Second)
	}

	var isRegistered bool
	for !isRegistered {
		time.Sleep(1 * time.Second)
		w.conn.Call("Coordinator.RegisterWorker", w.id, &isRegistered)
	}
}

func (w *Worker) WaitForInstructions() t.Instruction { // this really doesn't wait
	w.l.Println("waiting for instructions")
	var instruction *t.Instruction
	err := w.conn.Call("Coordinator.AskForInstructions", w.id, &instruction)
	if err != nil {
		w.l.Fatalf("error: %v", err)
	}
	return *instruction
}

/*
*
lifecycle of a worker

waits for instructions
executes the instructions
notifies the controller it is done
waits for instructions
*/
func (w *Worker) Loop() {
	// wait for instruction
	for {
		ins := w.WaitForInstructions()
		err := w.Execute(ins)
		if err != nil {
			w.Notify(false)
		}
		w.Notify(true)
	}
}

func (w *Worker) Execute(ins t.Instruction) error {
	w.l.Printf("I, worker with name %s, shall execute some tasks bestowed upon me by my master", w.id)
	return nil
}

func (w *Worker) Notify(success bool) {
	w.conn.Call("Coordinator.NotifyTaskStatus", t.Notification{w.id, success}, map[string]interface{}{})
}
