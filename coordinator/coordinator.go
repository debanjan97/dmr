package coordinator

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"

	t "github.com/debanjan97/distributed_map_reduce/types"
)

type InstructionSet struct {
	m            sync.Mutex
	Instructions []t.Instruction
}

type Coordinator struct {
	l                     log.Logger
	AvailableWorkers      []Worker
	InstructionSet        InstructionSet
	CompletedInstructions struct {
		m     sync.Mutex
		count int
	}
}

type Worker struct {
	id        string
	available bool
}

func (c *Coordinator) Serve() {
	rpc.Register(c)
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	go http.Serve(l, nil)
}

func (c *Coordinator) RegisterWorker(wid string, reply *bool) error {
	c.l.Printf("recieved registration request for worker %s", wid)
	c.AvailableWorkers = append(c.AvailableWorkers, Worker{wid, true})
	*reply = true
	return nil
}

func (c *Coordinator) AskForInstructions(wid string, reply *t.Instruction) error {
	// pop from the queue
	c.InstructionSet.m.Lock()
	defer c.InstructionSet.m.Unlock()
	if len(c.InstructionSet.Instructions) == 0 {
		return fmt.Errorf("no instructions")
	}
	*reply = c.InstructionSet.Instructions[0]
	c.InstructionSet.Instructions = c.InstructionSet.Instructions[1:] // pop

	for _, w := range c.AvailableWorkers {
		if w.id == wid && w.available {
			w.available = false
		}
	}

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

func (c *Coordinator) Start() {
	// For this working example, we need to print 10 times
	for i := 0; i < 10; i++ {
		c.InstructionSet.Instructions = append(c.InstructionSet.Instructions, t.Instruction{
			Operation: "PRINT",
			Operand:   []interface{}{""},
			Status:    "unclaimed",
		})
	}

	c.l.Print(len(c.InstructionSet.Instructions))
}

func (c *Coordinator) NotifyTaskStatus(n t.Notification, reply any) {
	for _, w := range c.AvailableWorkers {
		if w.id == n.WorkerId {
			w.available = true
		}
	}
	c.CompletedInstructions.m.Lock()
	defer c.CompletedInstructions.m.Lock()
	c.CompletedInstructions.count++
}

func (c *Coordinator) Done() bool {
	return c.CompletedInstructions.count == len(c.InstructionSet.Instructions)
}
