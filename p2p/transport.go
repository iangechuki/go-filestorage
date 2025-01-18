package p2p

// Peer is an interface that represents the remote node
type Peer interface {
}

// Transport is anything that handles the commication
// between nodes in the network.This can be in the form
// of TCP, UDP, WebRTC

type Transport interface {
	ListenAndAccept()
	Consume() <-chan RPC
}
