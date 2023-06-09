package p2p

import (
	"io"
	"net"
)

type Message struct {
	Payload io.Reader
	From    net.Addr
}

type Handshake struct {
	Version     string
	GameVariant GameVariant
	GameStatus  GameStatus
}
