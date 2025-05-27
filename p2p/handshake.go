package p2p

import "errors"

type HandshakeFunc func(Peer) error

func NopHandshake(Peer) error { return nil }

func Invalidhandshake(Peer) error {
	return errors.New("UnAuthorizes Access")
}
