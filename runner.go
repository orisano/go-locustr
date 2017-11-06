package locustr

import (
	"context"
	"net"
)

func Run(ctx context.Context, conn net.Conn) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	return nil
}
