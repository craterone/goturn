package libs

import "net"




type Allocate struct {
	Token string
	ClientAddress *net.UDPAddr
	PeerAddress *net.UDPAddr
	Peer *Peer
}

//func AllocatesMap2JSON(allocates map[string]*Allocate) []byte {
//
//}