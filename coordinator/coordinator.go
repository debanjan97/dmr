package coordinator

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
)

type Worker struct {
	id string
}

type Coordinator struct {
	l                log.Logger
	AvailableWorkers []Worker
}

func (c *Coordinator) Serve() {
	rpc.Register(c)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	http.Serve(l, nil)
}

func (c *Coordinator) RegisterWorker(wid string, reply *bool) error {
	c.l.Printf("recieved registration request for worker %s", wid)
	c.AvailableWorkers = append(c.AvailableWorkers, Worker{wid})
	*reply = true
	return nil
}

func (c *Coordinator) AskForInstructions(wid string, reply *map[string]interface{}) error {
	// not impl
	return nil
}

func NewCoordinator() *Coordinator {
	logger := log.Default()
	logger.SetPrefix("coordinator")
	return &Coordinator{
		l:                *log.New(logger.Writer(), logger.Prefix(), logger.Flags()),
		AvailableWorkers: make([]Worker, 0),
	}
}
