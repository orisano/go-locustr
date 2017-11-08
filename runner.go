package locustr

import (
	"context"
	"net"
)

func Run(ctx context.Context, conn net.Conn) error {
	client := &Client{
		rw:     conn,
		nodeID: GenNodeID(),
	}
	return client.Listen(ctx)
}

type runner struct {
	numClients int
	hatchRate  float64
}
