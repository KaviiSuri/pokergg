package p2p

import (
	"bytes"
	"encoding/gob"
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

	gameState *GameState
}

func NewServer(cfg ServerConfig) *Server {
	s := &Server{
		ServerConfig: cfg,
		peers:        make(map[net.Addr]*Peer),
		addPeer:      make(chan *Peer),
		delPeer:      make(chan *Peer),
		msgCh:        make(chan *Message),
		gameState:    NewGameState(),
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
		conn:     conn,
		outbound: true,
	}

	s.addPeer <- peer

	return s.SendHandshake(peer)
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

func (s *Server) SendHandshake(p *Peer) error {
	hs := &Handshake{
		GameVariant: s.GameVariant,
		Version:     s.Version,
		GameStatus:  s.gameState.gameStatus,
	}

	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(hs); err != nil {
		return err
	}

	return p.Send(buf.Bytes())
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
			// if a new player connects to the server, we send our handshake message
			// and wait for his reply
			if err := s.handleHandshake(peer); err != nil {
				logrus.Errorf("%s: handshake with incoming player failed: %s", s.ListenAddr, err)
				peer.conn.Close()
				delete(s.peers, peer.conn.RemoteAddr())
				continue
			}

			// TODO: Check max players and other game state logic
			go peer.ReadLoop(s.msgCh)

			if !peer.outbound {
				if err := s.SendHandshake(peer); err != nil {
					logrus.Errorf("failed to send handshake with peer: %s", err)
					peer.conn.Close()
					delete(s.peers, peer.conn.RemoteAddr())
					continue
				}
			}

			logrus.WithFields(logrus.Fields{
				"addr": peer.conn.RemoteAddr(),
			}).Info("new player connected")

			s.peers[peer.conn.RemoteAddr()] = peer
		case msg := <-s.msgCh:
			if err := s.HandleMessage(msg); err != nil {
				panic(err)
			}
		}
	}
}

func (s *Server) handleHandshake(p *Peer) error {
	hs := &Handshake{}
	if err := gob.NewDecoder(p.conn).Decode(hs); err != nil {
		return err
	}

	if s.GameVariant != hs.GameVariant {
		return fmt.Errorf("game variant mismatch server: %s client: %s", s.GameVariant, hs.GameVariant)
	}

	// TODO: Add sematic versioning support
	if s.Version != hs.Version {
		return fmt.Errorf("version mismatch server: %s client: %s", s.Version, hs.Version)
	}

	logrus.WithFields(logrus.Fields{
		"peer":       p.conn.RemoteAddr(),
		"version":    hs.Version,
		"variant":    hs.GameVariant,
		"gameStatus": hs.GameStatus,
	}).Info("handshake recieved")

	return nil
}

func (s *Server) HandleMessage(msg *Message) error {
	fmt.Printf("%+v\n", msg)
	return nil
}
