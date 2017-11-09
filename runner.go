package locustr

import (
	"context"
	"math"
	"net"
	"sync"

	"github.com/pkg/errors"
)

type State int

const (
	StateInit State = iota
	StateHatching
	StateRunning
	StateStopped
)

func Run(ctx context.Context, conn net.Conn) error {
	client := &Client{
		rw:     conn,
		nodeID: GenNodeID(),
	}
	return client.Listen(ctx)
}

type Runner struct {
	conn  net.Conn
	tasks []Task

	nodeID string

	hatchRate float64

	state   State
	stateMu sync.RWMutex

	locusts   []Locust
	locustsMu sync.RWMutex
}

func NewRunner(masterHost string, masterPort int) (*Runner, error) {
	conn, err := NewSocket(masterHost, masterPort)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect master")
	}
	nodeID := GenNodeID()

	return &Runner{
		conn:      conn,
		nodeID:    nodeID,
		hatchRate: 1,
		state:     StateInit,
	}, nil
}

func (r *Runner) weightLocusts(amount int) []Locust {
	sum := 0.0
	for _, task := range r.tasks {
		sum += task.Weight
	}

	locusts := make([]Locust, 0, amount+10)
	for _, task := range r.tasks {
		percent := task.Weight / sum
		n := int(math.Floor(float64(amount)*percent + 0.5))
		locust := Locust{
			task: task,
		}
		for i := 0; i < n; i++ {
			locusts = append(locusts, locust)
		}
	}
	return locusts
}

func (r *Runner) spawnLocusts(count int) {
	bucket := r.weightLocusts(count)
	count = len(bucket)
	if r.state == StateInit || r.state == StateStopped {
		r.state = StateHatching
	}
}
