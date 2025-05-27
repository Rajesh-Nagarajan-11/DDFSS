package p2p

// Peer is an interface that represents any remote node
type Peer interface {
	Close() error
}

// Transport is anything that handles communication between nodes in the network
// This can be a from of (TCP,UDP,Websocket)
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
}
