package libs

import (
	"net"
	"strconv"
	"log"
)

//relay
type Relay struct {
	Port       int
	connection *net.UDPConn
	RelayAddress *net.UDPAddr
	ServerAddress *net.UDPAddr
}

func (s *Relay) serve() {
	for {
		var buf = make([]byte, 2048)
		size, remoteAddr, err := s.connection.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		go s.handleData(remoteAddr, buf[:size])
	}
}

//Serve initiates a UDP connection that listens on any port for incoming data
func (s *Relay) Serve() {
	laddr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(s.Port))
	if err != nil {
		log.Fatal(err)
	}
	s.connection, err = net.ListenUDP("udp", laddr)
	if err != nil {
		log.Fatal(err)
	}
	go s.serve()
}

func (s *Relay) handleData(raddr *net.UDPAddr, data []byte) {

	switch data[0] {
	//from client or 3478
	case 0x00:
		msg , err := UnMarshal(data)
		if err != nil {
			Log.Warning(err)
			return
		}
		switch msg.MessageType {
		case TypeBindingRequest:
			//todo : ignore?
		case TypeSendIndication:
			//Log.Infof("peer port %d ---> request : %s \n",s.Port,msg)
			data = append([]byte{0x11},data...)

			if s.RelayAddress != nil{
				_, err := s.connection.WriteToUDP(data, s.RelayAddress)
				if err != nil {
					Log.Warning(err)
				}
			}else{
				//fixme firefox drop packet
				Log.Infof("ffffff")
			}


		}

	case 0x11:
		if s.ServerAddress != nil{
			_, err := s.connection.WriteToUDP(data, s.ServerAddress)
			if err != nil {
				Log.Warning(err)
			}
		}else{
			Log.Infof("ffffff")
		}
	case 0x40:
		data = append([]byte{0x11},data...)

		if s.RelayAddress != nil{
			_, err := s.connection.WriteToUDP(data, s.RelayAddress)
			if err != nil {
				Log.Warning(err)
			}
		}else{
			Log.Infof("ffffff")
		}

	}



}