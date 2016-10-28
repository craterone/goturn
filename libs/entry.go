package libs


import (
	"log"
	"net"
	"strconv"
)



type Entry struct {
	Port       int
	udpConn *net.UDPConn
}

func LoadEntryModule()  {

	entry := NewEntry(*server_port)
	entry.Serve()
}

func (s *Entry) serveUDP() {
	laddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(s.Port))
	if err != nil {
		log.Fatal(err)
	}
	s.udpConn, err = net.ListenUDP("udp", laddr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		var buf = make([]byte, 2048)
		size, remoteAddr, err := s.udpConn.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		go s.handleData(remoteAddr, buf[:size],false)
	}
}

func (s *Entry) Serve() {
	serverTCP, serverTLS := false,false
	serverUDP := true

	if serverTCP{

	}

	if serverTLS {

	}

	if serverUDP {
		s.serveUDP()
	}

}

func NewEntry(port int) *Entry {
	ret := new(Entry)
	ret.Port = port
	return ret
}

func (entry *Entry) handleData(raddr *net.UDPAddr, data []byte,tcp bool) {
	//check packet from
	switch data[0] {
	//from client
	case 0x00:
		msg, err := UnMarshal(data)
		if err != nil {
			Log.Warning(err)
			return
		}


		var response []byte
		var responseError error
		var responseAddress *net.UDPAddr

		switch msg.MessageType {
		case TypeBindingRequest:
			//todo : ignore
		case TypeAllocate , TypeCreatePermission , TypeChannelBinding, TypeRefreshRequest:
			response,responseError = turnMessageHandle(msg,raddr,false)
		case TypeSendIndication:
			response,responseAddress,responseError = turnRelayMessageHandle(msg,raddr,false)

		}

		if responseError == nil{
			if !tcp {
				if response != nil {
					if responseAddress != nil{
						raddr = responseAddress
					}
					_, err := entry.udpConn.WriteToUDP(response, raddr)
					if err != nil {
						Log.Warning(err)
					}
				}else {
					//todo add message type check
					Log.Warningf("no response.  with %s",msg.TypeToString())
				}

			}else {
				//todo : add tcp
			}
		}else{
			Log.Warningf("response error : %v",responseError)
		}


	//from other peers
	case 0x11:
		extractData := data[1:]

		if extractData[0] == 0x40 {
			for k,v := range GlobalAllocates {
				if v.PeerAddress.String() == raddr.String() {
					responseAddr := parseAddressV4(k)

					_, err := entry.udpConn.WriteToUDP(extractData, responseAddr)
					if err != nil {
						Log.Warning(err)
					}
					break
				}
			}
		}else{
			msg, err := UnMarshal(extractData)
			if err != nil {
				Log.Warning(err)
				return
			}

			for k,v := range GlobalAllocates {
				if v.PeerAddress.String() == raddr.String() {
					responseAddr := parseAddressV4(k)

					respMessage := NewResponse(TypeDataIndication,msg.TransID,
						msg.getAttribute(AttributeData),
						msg.getAttribute(AttributeXorPeerAddress),
						AttrSoftware,
					)

					response, _ := Marshal(respMessage,false)

					_, err := entry.udpConn.WriteToUDP(response, responseAddr)
					if err != nil {
						Log.Warning(err)
					}
					break
				}
			}
		}

	//from channel
	case 0x40:
		clientAddress := raddr.String()

		allocate := GlobalAllocates[clientAddress]

		if allocate != nil {
			_, err := entry.udpConn.WriteToUDP(data, allocate.PeerAddress)
			if err != nil {
				Log.Warning(err)
			}
		}
	}

}
