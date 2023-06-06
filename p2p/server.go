package p2p

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
)

type GameVariant uint8

const (
	TexasHoldem GameVariant = iota
	Other
)

func (gv GameVariant) String() string {
	switch gv {
	case TexasHoldem:
		return "Texas Hold'em"
	default:
		return "Unknown"
	}
}

type ServerConfig struct {
	ListenAddr  string
	Version     string
	GameVariant GameVariant
}

type Server struct {
	ServerConfig

	transport *TCPTransport
	peers     map[net.Addr]*Peer
	addPeer   chan *Peer
	delPeer   chan *Peer
	msgCh     chan *Message
}

func NewServer(cfg ServerConfig) *Server {
	s := &Server{
		ServerConfig: cfg,
		peers:        make(map[net.Addr]*Peer),
		addPeer:      make(chan *Peer),
		delPeer:      make(chan *Peer),
		msgCh:        make(chan *Message),
	}

	tr := NewTCPTransport(s.ListenAddr)
	s.transport = tr

	tr.AddPeer = s.addPeer
	tr.DelPeer = s.delPeer

	return s
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

	err = peer.Send([]byte(s.Version))
	if err != nil {
		panic(err)
	}
	return err
}

func (s *Server) Start() {
	go s.loop()

	logrus.WithFields(logrus.Fields{
		"addr": s.ListenAddr,
		"type": s.GameVariant,
	}).Info("started new game server")
	err := s.transport.ListenAndAccept()
	if err != nil {
		panic(err)
	}
}

func (s *Server) loop() {

	for {
		select {
		case peer := <-s.delPeer:
			addr := peer.conn.RemoteAddr()
			logrus.WithFields(logrus.Fields{
				"addr": addr,
			}).Info("player disconnected")
			delete(s.peers, addr)
		case peer := <-s.addPeer:
			// TODO: Check max players and other game state logic
			logrus.WithFields(logrus.Fields{
				"addr": peer.conn.RemoteAddr(),
			}).Info("new player connected")
			go peer.ReadLoop(s.msgCh)
			s.peers[peer.conn.RemoteAddr()] = peer
		case msg := <-s.msgCh:
			if err := s.HandleMessage(msg); err != nil {
				panic(err)
			}
		}
	}
}

func (s *Server) HandleMessage(msg *Message) error {
	fmt.Printf("%+v\n", msg)
	return nil
}
