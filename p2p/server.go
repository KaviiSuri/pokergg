package p2p

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

type Peer struct {
	conn net.Conn
}

func (p *Peer) Send(data []byte) error {
	_, err := p.conn.Write(data)
	return err
}

type ServerConfig struct {
	ListenAddr string
	Version    string
}

type Message struct {
	Payload io.Reader
	From    net.Addr
}

type Server struct {
	ServerConfig

	handler  Handler
	listener net.Listener
	peers    map[net.Addr]*Peer
	addPeer  chan *Peer
	delPeer  chan *Peer
	msgCh    chan *Message
}

func NewServer(cfg ServerConfig) *Server {
	return &Server{
		ServerConfig: cfg,
		peers:        make(map[net.Addr]*Peer),
		addPeer:      make(chan *Peer),
		handler:      &DefaultHandler{},
		msgCh:        make(chan *Message),
	}
}

// TODO: redundant code in Connect and acceptLoop
func (s *Server) Connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	peer := &Peer{
		conn: conn,
	}

	s.addPeer <- peer

	return peer.Send([]byte(fmt.Sprintf("%s\n", s.Version)))
}

func (s *Server) Start() {
	go s.loop()

	if err := s.listen(); err != nil {
		panic(err)
	}

	fmt.Printf("game server running on %s\n", s.ListenAddr)

	s.acceptLoop()
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			panic(err)
		}

		peer := &Peer{
			conn: conn,
		}

		s.addPeer <- peer

		peer.Send([]byte(fmt.Sprintf("%s\n", s.Version)))

		go s.handleConn(peer)
	}
}

func (s *Server) handleConn(peer *Peer) {
	buf := make([]byte, 1024)
	for {
		n, err := peer.conn.Read(buf)
		if err != nil {
			break
		}

		s.msgCh <- &Message{
			Payload: bytes.NewReader(buf[:n]),
			From:    peer.conn.RemoteAddr(),
		}
	}

	s.delPeer <- peer
}

func (s *Server) listen() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}

	s.listener = ln

	return nil
}

func (s *Server) loop() {

	for {
		select {
		case peer := <-s.delPeer:
			addr := peer.conn.RemoteAddr()
			delete(s.peers, addr)
			fmt.Printf("player disconnected %s\n", addr)
		case peer := <-s.addPeer:
			fmt.Printf("new player connected %s\n", peer.conn.RemoteAddr())
			s.peers[peer.conn.RemoteAddr()] = peer
		case msg := <-s.msgCh:
			if err := s.handler.HandleMessage(msg); err != nil {
				panic(err)
			}
		}
	}
}
