package locustr

import (
	"io"

	"github.com/ugorji/go/codec"
)

var messageMsgpackHandle codec.MsgpackHandle

type Message struct {
	_struct bool                   `codec:",toarray"`
	Type    string                 `codec:"type"`
	Data    map[string]interface{} `codec:"data"`
	NodeID  string                 `codec:"node_id"`
}

func EncodeMessage(w io.Writer, m *Message) error {
	return codec.NewEncoder(w, &messageMsgpackHandle).Encode(m)
}

func DecodeMessage(r io.Reader, m *Message) error {
	return codec.NewDecoder(r, &messageMsgpackHandle).Decode(m)
}
