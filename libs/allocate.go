package libs

import "net"




type Allocate struct {
	ClientAddress *net.UDPAddr
	PeerAddress *net.UDPAddr
	Peer *Peer
}

