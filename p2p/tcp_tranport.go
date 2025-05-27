package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TcpPeer Represents Remote Node of the TCP Establisted Connection
type TcpPeer struct {
	// This field represents the actual TCP connection between your device and the remote device (peer).
	// Think of it as the phone line between your device and the other device.
	conn net.Conn
	// if we dial and retrive connection => outbound:true
	// if we accept and retrive connection =>outbound false
	// tells whether the connection was initiated by you (outbound) or received by you (inbound) i.e outbound => false.
	outbound bool
}

func NewTcpPeer(conn net.Conn, outbound bool) *TcpPeer {
	return &TcpPeer{
		conn:     conn,
		outbound: outbound,
	}
}

// close implements peer interface
func (P *TcpPeer) Close() error {
	return P.conn.Close()
}

type TCP_Transportopts struct {
	// The TCP address this transport will bind to and listen for incoming peer connections (e.g., ":3000").
	Listenaddr string
	// A function that defines the handshake logic when a peer connects (used for authentication, version checking, etc.).
	Handshakefunc HandshakeFunc
	// Responsible for decoding the incoming messages from peers into a format your application understands.
	Decoder Decoder
	// A callback that is triggered after a successful handshake with a peer; used to manage peer behavior.
	Onpeer func(Peer) error
}

type TCP_Transport struct {
	TCP_Transportopts                   // Embedded configuration options for the transport
	listener          net.Listener      // The network listener to accept incoming connections
	rpcch             chan RPC          // Channel used to send or receive RPC (Remote Procedure Call) messages between peers
	mu                sync.RWMutex      // A lock to protect shared data (used for synchronization)
	peers             map[net.Addr]Peer // A map to store connected peers (other devices or systems)
}

func NewTCP_Transport(opts TCP_Transportopts) *TCP_Transport {

	return &TCP_Transport{
		TCP_Transportopts: opts,
		rpcch:             make(chan RPC),
	}
}

// consume implements transport interface, which will return a read only channel
// for reading incoming messages received from other peer from network
func (t *TCP_Transport) Consume() <-chan RPC {
	return t.rpcch
}

func (t *TCP_Transport) ListenAndAccept() error {
	var err error
	t.listener, err = net.Listen("tcp", t.Listenaddr)
	if err != nil {
		fmt.Printf("Failed to listen on address %s: %v\n", t.Listenaddr, err)
		return err
	}
	go t.startAcceptLoop()
	return nil
}

func (t *TCP_Transport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP Accept Error :%s\n", err)
		}
		fmt.Printf("New Incoming Connection %+v\n", conn)
		go t.handleConn(conn)
	}
}

func (t *TCP_Transport) handleConn(conn net.Conn) {
	var err error

	defer func() {
		fmt.Printf("Dropping Peer Connection %s\n", err)
		conn.Close()
	}()

	peer := NewTcpPeer(conn, true)

	if err = t.Handshakefunc(peer); err != nil {

		return
	}

	if t.Onpeer != nil {
		if err = t.Onpeer(peer); err != nil {
			return
		}
	}

	rpc := RPC{}
	for {

		if err := t.Decoder.Decode(conn, &rpc); err != nil {
			fmt.Printf("TCP Error %s\n", err)
			continue
		}
		rpc.From = conn.RemoteAddr()
		t.rpcch <- rpc
	}

}
