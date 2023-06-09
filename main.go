package main

import (
	"fmt"
	"time"

	"github.com/KaviiSuri/pokergg/p2p"
)

func main() {
	cfg := p2p.ServerConfig{
		ListenAddr:  ":3000",
		Version:     "POKERGG V0.1-alpha",
		GameVariant: p2p.TexasHoldem,
	}
	server := p2p.NewServer(cfg)
	go server.Start()

	time.Sleep(1 * time.Second)

	remoteCfg := p2p.ServerConfig{
		ListenAddr:  ":4000",
		Version:     "POKERGG V0.1-alpha",
		GameVariant: p2p.TexasHoldem,
	}
	remoteServer := p2p.NewServer(remoteCfg)
	go remoteServer.Start()
	if err := remoteServer.Connect("localhost:3000"); err != nil {
		fmt.Println(err)
	}

	select {}
	//rand.Seed(time.Now().UnixNano())
	//for i := 0; i < 10; i++ {
	//	d := deck.New()
	//	fmt.Println(d)
	//}
}
