package locustr

import (
	"context"
	"math"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/orisano/worker"
	"github.com/pkg/errors"
)

type State int

const (
	StateInit State = iota
	StateHatching
	StateRunning
	StateStopped
)

type Runner struct {
	conn  net.Conn
	tasks []Task

	nodeID string

	numClients int32
	hatchRate  float64

	state   State
	stateMu sync.RWMutex

	locusts    worker.Group
	occurrence map[string][]int32
}

func NewRunner(masterHost string, masterPort int) (*Runner, error) {
	conn, err := NewSocket(masterHost, masterPort)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect master")
	}
	nodeID := GenNodeID()

	// context null check only
	locusts, _ := worker.NewGroup(context.Background())

	return &Runner{
		conn:       conn,
		nodeID:     nodeID,
		numClients: 0,
		hatchRate:  1,
		state:      StateInit,
		locusts:    locusts,
		occurrence: make(map[string][]int32),
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
	shuffle(bucket)
	if r.state == StateInit || r.state == StateStopped {
		r.state = StateHatching
		r.numClients = int32(count)
	} else {
		r.numClients += int32(count)
	}

	interval := time.Duration(int64(float64(time.Second) / r.hatchRate))
	ticker := time.NewTicker(interval)
	for _, locust := range bucket {
		id := r.locusts.Spawn(locust.Run)
		r.occurrence[locust.task.Name] = append(r.occurrence[locust.task.Name], id)
		<-ticker.C
	}
	// sendHatchCompleteEvent r.numClients
}

func shuffle(locusts []Locust) {
	for i := len(locusts) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		locusts[i], locusts[j] = locusts[j], locusts[i]
	}
}
