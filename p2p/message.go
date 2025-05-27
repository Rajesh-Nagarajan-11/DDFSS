package p2p

import "net"

// RPC holds any arbitrary data that is being sents over the
// each transport between two nodes in networks
type RPC struct {
	From    net.Addr
	Payload []byte
}
