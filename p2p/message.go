package p2p

import "net"

// Message holds any arbitrary data that is being sent over the network
// between peers
type RPC struct {
	From    net.Addr
	Payload []byte
}
