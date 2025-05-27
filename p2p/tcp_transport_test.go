package p2p

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTcpTransport(t *testing.T) {
	opts := TCP_Transportopts{
		Listenaddr:    ":3000",
		Handshakefunc: NopHandshake,
		Decoder:       DefaultDecoder{},
	}

	transport := NewTCP_Transport(opts)

	assert.Equal(t, ":3000", transport.Listenaddr)

	go func() {
		err := transport.ListenAndAccept()
		assert.Nil(t, err, "Expected no error but got %v", err)
	}()

	// Allow some time for the listener to run (e.g., simulate a running server)
	time.Sleep(1 * time.Second)
}
