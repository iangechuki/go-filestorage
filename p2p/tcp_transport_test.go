package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	listenAddr := "4008"
	tr := NewTCPTransport(listenAddr)

	assert.Equal(t, listenAddr, tr.listenAddress)

	// server
	// tr.Start()
	assert.Nil(t, tr.ListenAndAccept())
}
