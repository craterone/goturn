package libs

import (
	"net"
	"encoding/hex"
	"bytes"
)

func messageIntegrityCheck(requestMessage *Message) (err error) {
	miAttr := requestMessage.getAttribute(AttributeMessageIntegrity)

	if miAttr != nil {
		userAttr := requestMessage.getAttribute(AttributeUsername)
		if userAttr != nil {
			username := string(userAttr.Value)
			password := hex.EncodeToString(HmacSha1(userAttr.Value,[]byte("passwordkey")))

			//Log.Infof("password %s",password)
			key := generateKey(username,password,"realm")
			requestValue , err :=  Marshal(requestMessage,true)

			//Log.Infof("password %s , username %s ",password,username)

			if err == nil {
				calculateMi := MessageIntegrityHmac(requestValue,key)

				//Log.Infof("origin %x , after %x ",miAttr.Value,calculateMi)

				if(!bytes.Equal(calculateMi,miAttr.Value)){
					//todo : not equal
				}
			}else{
				//todo
			}
		}else{
			//todo
		}
	}else{
		//todo : error response
	}

	return
}

func messageIntegrityCalculate(username string,responseMessage *Message) (response []byte, err error) {
	var m_i_response []byte
	m_i_response, err = Marshal(responseMessage,false)

	if err != nil {
		return nil,err
	}

	password := hex.EncodeToString(HmacSha1([]byte(username),[]byte("passwordkey")))

	key := generateKey(username,password,"realm")

	hmacValue := MessageIntegrityHmac(m_i_response[:len(m_i_response)-24],key)

	response = append(m_i_response[:len(m_i_response)-20],hmacValue...)

	return response,nil
}

var (
	xxxxx = 0;
	xxxxaaa = []int{11111,11112,11113}
)
func turnMessageHandle(requestMessage *Message,clientAddr *net.UDPAddr,tcp bool) ([]byte, error) {

	switch requestMessage.MessageType {
	case TypeAllocate:
		//long-term , check username
		usernameAttr := requestMessage.getAttribute(AttributeUsername)

		if usernameAttr != nil {

			//todo : add check
			usernameStr := string(usernameAttr.Value)


			relayPort := xxxxaaa[xxxxx]
			xxxxx++
			//relayPort := RelayPortPool.RandSelectPort()
			relayAddress := getRelayAddress()

			peerAddress := new(net.UDPAddr)
			peerAddress.Port = relayPort
			peerAddress.IP = net.ParseIP(relayAddress)


			peer := new(Peer)
			peer.Port = relayPort
			peer.ServerAddress = ServerAddress
			peer.Serve()




			allocate := new(Allocate)
			allocate.ClientAddress = clientAddr
			allocate.Peer = peer
			allocate.PeerAddress = peerAddress

			clientAddrStr := clientAddr.String()

			AllocateMutex.Lock()
			GlobalAllocates[clientAddrStr] = allocate
			AllocateMutex.Unlock()

			respMsg := NewResponse(TypeAllocateResponse,requestMessage.TransID,
				newAttrXORRelayedAddress(relayAddress,relayPort),
				newAttrXORMappedAddress(clientAddr),
				AttrLifetime,
				AttrSoftware,
				AttrDummyMessageIntegrity,
			)


			response,err := messageIntegrityCalculate(usernameStr,respMsg)

			return response,err

		}else{
			// 401 error
			respMsg := NewResponse(TypeAllocateErrorResponse,requestMessage.TransID,
				newAttrNonce(),
				AttrRealm,
				AttrError401,
				AttrSoftware,
			)

			response, err := Marshal(respMsg,false)

			if err != nil {
				return nil,err
			}
			return response,nil
		}
	case TypeCreatePermission:
		respMsg := NewResponse(TypeCreatePermissionResponse,requestMessage.TransID,
			AttrSoftware,
			AttrDummyMessageIntegrity,
		)

		originUsername := requestMessage.getAttribute(AttributeUsername)
		strUsername := string(originUsername.Value)

		response ,err := messageIntegrityCalculate(strUsername,respMsg)
		return response,err
	case TypeChannelBinding:

		respMsg := NewResponse(TypeChannelBindingResponse,requestMessage.TransID,
			AttrSoftware,
			AttrDummyMessageIntegrity,
		)

		originUsername := requestMessage.getAttribute(AttributeUsername)
		strUsername := string(originUsername.Value)

		response ,err := messageIntegrityCalculate(strUsername,respMsg)
		return response,err
	case TypeRefreshRequest:

		respMsg := NewResponse(TypeRefreshResponse,requestMessage.TransID,
			AttrSoftware,
			newAttrLifetime(),
			AttrDummyMessageIntegrity,
		)

		originUsername := requestMessage.getAttribute(AttributeUsername)
		strUsername := string(originUsername.Value)

		response ,err := messageIntegrityCalculate(strUsername,respMsg)
		return response,err

	}

	return nil,nil
}

func turnRelayMessageHandle(requestMessage *Message,clientAddr *net.UDPAddr,tcp bool) ([]byte, *net.UDPAddr,error) {

	switch requestMessage.MessageType {
	case TypeSendIndication:
		clientAddress := clientAddr.String()

		allocate := GlobalAllocates[clientAddress]

		if allocate.Peer != nil{
			if allocate.Peer.RelayAddress == nil {
				peerAddress := requestMessage.getAttribute(AttributeXorPeerAddress)

				if peerAddress != nil {
					port, address := unXorAddress(peerAddress.Value)
					//relayAddress := fmt.Sprintf("%s:%d",net.IP(address),port)

					relayAddress := new(net.UDPAddr)
					relayAddress.Port = int(port)
					relayAddress.IP = address

					allocate.Peer.RelayAddress = relayAddress


					Log.Infof("relay address : %s , peer addres : %s",relayAddress,allocate.PeerAddress.String())
				}

			}

			requestMessage.setAttribute(newAttrXORPeerAddress(allocate.PeerAddress.IP.String(),allocate.PeerAddress.Port))

			response, err := Marshal(requestMessage,false)

			return response,allocate.PeerAddress,err
		}
	}
	return nil,nil,nil
}

