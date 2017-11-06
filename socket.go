package locustr

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
)

func NewSocket(host string, port int) (net.Conn, error) {
	saddr := fmt.Sprintf("%s:%d", host, port)
	addr, err := net.ResolveTCPAddr("tcp", saddr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to resolve addr")
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect")
	}
	conn.SetNoDelay(true)
	return conn, nil
}
