package locustr

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Environment struct {
	HostName func() (string, error)
	TimeNow  func() time.Time
	RandInt  func(n int64) int64
}

var DefaultEnvironment = &Environment{
	HostName: os.Hostname,
	TimeNow:  time.Now,
	RandInt:  rand.Int63n,
}

// ref: https://github.com/locustio/locust/blob/master/locust/runners.py#L355
func (e *Environment) GenNodeID() string {
	hostname, _ := e.HostName()
	timestamp := e.TimeNow().Unix()
	randomNum := e.RandInt(10000)

	b := new(bytes.Buffer)
	b.WriteString(hostname)
	b.WriteByte('_')

	d := md5.New()
	io.WriteString(d, strconv.FormatInt(timestamp+randomNum, 10))
	digest := d.Sum(nil)
	hexDigest := hex.EncodeToString(digest)

	b.WriteString(hexDigest)
	return b.String()
}

func GenNodeID() string {
	return DefaultEnvironment.GenNodeID()
}
