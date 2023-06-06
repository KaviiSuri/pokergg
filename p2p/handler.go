package p2p

import (
	"fmt"
	"io"
)

type Handler interface {
	HandleMessage(*Message) error
}

type DefaultHandler struct {
}

func (h *DefaultHandler) HandleMessage(m *Message) error {
	b, err := io.ReadAll(m.Payload)
	if err != nil {
		return err
	}
	fmt.Printf("received message from %s: %s\n", m.From, string(b))
	return nil
}
