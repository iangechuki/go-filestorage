package p2p

// Message holds any arbitrary data that is being sent over the network
// between peers
type RPC struct {
	From    string
	Payload []byte
}
