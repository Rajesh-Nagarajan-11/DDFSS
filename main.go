package main

import (
	"fmt"
	"log"

	"ddfss.com/p2pTcp/p2p"
)

func Onpeer(peer p2p.Peer) error {
	//peer.Close()
	return nil
}

func main() {
	tcpopts := p2p.TCP_Transportopts{
		Listenaddr:    ":3000",
		Handshakefunc: p2p.NopHandshake,
		Decoder:       p2p.DefaultDecoder{},
		Onpeer:        Onpeer,
	}
	tr := p2p.NewTCP_Transport(tcpopts)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("msg : %+v\n", msg)
		}
	}()
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	select {}
}
