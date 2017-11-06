package locustr

import (
	"bufio"
	"context"
	"encoding/binary"
	"io"
	"io/ioutil"
	"net"

	"github.com/pkg/errors"
)

type Client struct {
	conn   net.Conn
	nodeID string
}

func (c *Client) Listen(ctx context.Context) error {
	var cancel context.CancelFunc
	safeCancel := func() {
		if cancel != nil {
			cancel()
		}
	}
	defer safeCancel()

	c.sendReady()

	br := bufio.NewReader(c.conn)
	buf := make([]byte, 4)
	var msg Message
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if _, err := io.ReadFull(br, buf); err != nil {
				return errors.Wrap(err, "failed to read connection")
			}
			length := binary.BigEndian.Uint32(buf)
			r := io.LimitReader(br, int64(length))
			if err := DecodeMessage(r, &msg); err == nil {
				switch msg.Type {
				case "hatch":
					workers, hatchRate := parseHatchData(&msg)
					var wctx context.Context
					wctx, cancel = context.WithCancel(ctx)
					go c.startHatching(wctx, workers, hatchRate)
				case "stop":
					safeCancel()
					c.sendStopped()
					c.sendReady()
				case "quit":
					return nil
				}
			}
			io.Copy(ioutil.Discard, r)
		}
	}
}

func (c *Client) startHatching(ctx context.Context, locustCount int64, hatchRate float64) error {
	return nil
}

func (c *Client) send(messageType string, data map[string]interface{}) error {
	return EncodeMessage(c.conn, &Message{
		Type:   messageType,
		Data:   data,
		NodeID: c.nodeID,
	})
}

func (c *Client) sendReady() error {
	return c.send("client_ready", nil)
}

func (c *Client) sendStopped() error {
	return c.send("client_stopped", nil)
}

func parseHatchData(m *Message) (int64, float64) {
	hatchRate := m.Data["hatch_rate"].(float64)
	clients := m.Data["num_clients"]
	if i, ok := clients.(int64); ok {
		return i, hatchRate
	} else {
		return int64(clients.(uint64)), hatchRate
	}
}
