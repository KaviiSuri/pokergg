package p2p

import (
	"bytes"
	"net"

	"github.com/sirupsen/logrus"
)

type Peer struct {
	conn     net.Conn
	outbound bool
}

func (p *Peer) Send(data []byte) error {
	_, err := p.conn.Write(data)
	return err
}

func (p *Peer) ReadLoop(msgCh chan *Message) {
	buf := make([]byte, 1024)
	for {
		n, err := p.conn.Read(buf)
		if err != nil {
			break
		}
		msg := &Message{
			Payload: bytes.NewReader(buf[:n]),
			From:    p.conn.RemoteAddr(),
		}
		msgCh <- msg
	}

	// TODO: unregister this peer
	p.conn.Close()
}

type TCPTransport struct {
	listenAddr string
	listener   net.Listener
	AddPeer    chan *Peer
	DelPeer    chan *Peer
}

func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		listenAddr: listenAddr,
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	ln, err := net.Listen("tcp", t.listenAddr)
	if err != nil {
		return err
	}
	t.listener = ln

	for {
		conn, err := ln.Accept()
		if err != nil {
			logrus.Error(err)
			continue
		}

		peer := &Peer{
			conn:     conn,
			outbound: false,
		}
		t.AddPeer <- peer
	}
}
