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

				Log.Infof("origin %x , after %x ",miAttr.Value,calculateMi)

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

func messageIntegrityCalculate(username []byte,responseMessage *Message) (response []byte, err error) {
	var m_i_response []byte
	m_i_response, err = Marshal(responseMessage,false)

	if err != nil {
		return nil,err
	}

	password := hex.EncodeToString(HmacSha1(username,[]byte("passwordkey")))

	key := generateKey(bytes2str(username),password,"realm")

	hmacValue := MessageIntegrityHmac(m_i_response[:len(m_i_response)-24],key)

	response = append(m_i_response[:len(m_i_response)-20],hmacValue...)

	return response,nil
}

func turnMessageHandle(requestMessage *Message,clientAddr *net.UDPAddr,tcp bool) ([]byte, error) {

	switch requestMessage.MessageType {
	case TypeAllocate:
		//long-term , check username
		uAttr := requestMessage.getAttribute(AttributeUsername)

		if uAttr != nil {
			//check username format
			icolon := bytes.IndexByte(uAttr.Value,':')

			if icolon < 0 {
				//todo : error
			}

			timestamp := strBytes2Int64(uAttr.Value[:icolon])
			token := uAttr.Value[icolon+1:]


			Log.Infof("timestamp : %d , username : %s",timestamp,token)
			//fixme : check time expire
			if timestamp < 0{
				//todo error
			}

			err := messageIntegrityCheck(requestMessage)
			if err != nil {
				//todo error
				Log.Infof("message integrity error")
			}

			rport := RelayPortPool.RandSelectPort()
			rip := getRelayAddress()

			raddress := new(net.UDPAddr)
			raddress.Port = rport
			raddress.IP = rip

			peer := new(Peer)
			peer.Port = rport
			peer.ServerAddress = ServerAddress
			peer.Serve()

			allocate := new(Allocate)
			allocate.Token = bytes2str(token)
			allocate.ClientAddress = clientAddr
			allocate.Peer = peer
			allocate.PeerAddress = raddress

			clientAddrStr := clientAddr.String()

			AllocateMutex.Lock()
			GlobalAllocates[clientAddrStr] = allocate
			AllocateMutex.Unlock()

			respMsg := NewResponse(TypeAllocateResponse,requestMessage.TransID,
				newAttrXORRelayedAddress(rip,rport),
				newAttrXORMappedAddress(clientAddr.IP.To4(),clientAddr.Port),
				AttrLifetime,
				AttrSoftware,
				AttrDummyMessageIntegrity,
			)


			response,err := messageIntegrityCalculate(uAttr.Value,respMsg)

			return response,err

		}else{
			// 401 error
			respMsg := NewResponse(TypeAllocateErrorResponse,requestMessage.TransID,
				AttrNonce,
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

		clientAddress := clientAddr.String()

		allocate := GlobalAllocates[clientAddress]

		err := messageIntegrityCheck(requestMessage)

		if err != nil {
			Log.Infof("message integrity error")
		}

		if allocate != nil {
			if allocate.Peer != nil {
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
		}

		respMsg := NewResponse(TypeCreatePermissionResponse,requestMessage.TransID,
			AttrSoftware,
			AttrDummyMessageIntegrity,
		)

		originUsername := requestMessage.getAttribute(AttributeUsername)

		response ,err := messageIntegrityCalculate(originUsername.Value,respMsg)
		return response,err
	case TypeChannelBinding:

		err := messageIntegrityCheck(requestMessage)

		if err != nil {
			Log.Infof("message integrity error")
		}

		respMsg := NewResponse(TypeChannelBindingResponse,requestMessage.TransID,
			AttrSoftware,
			AttrDummyMessageIntegrity,
		)

		originUsername := requestMessage.getAttribute(AttributeUsername)

		response ,err := messageIntegrityCalculate(originUsername.Value,respMsg)
		return response,err
	case TypeRefreshRequest:

		err := messageIntegrityCheck(requestMessage)

		if err != nil {
			Log.Infof("message integrity error")
		}

		respMsg := NewResponse(TypeRefreshResponse,requestMessage.TransID,
			AttrSoftware,
			newAttrLifetime(),
			AttrDummyMessageIntegrity,
		)

		originUsername := requestMessage.getAttribute(AttributeUsername)

		response ,err := messageIntegrityCalculate(originUsername.Value,respMsg)
		return response,err

	}

	return nil,nil
}

func turnRelayMessageHandle(requestMessage *Message,clientAddr *net.UDPAddr,tcp bool) ([]byte, *net.UDPAddr,error) {

	switch requestMessage.MessageType {
	case TypeSendIndication:
		//todo : check permission

		clientAddress := clientAddr.String()

		allocate := GlobalAllocates[clientAddress]

		if allocate.Peer != nil{

			requestMessage.setAttribute(newAttrXORPeerAddress(allocate.PeerAddress.IP.To4(),allocate.PeerAddress.Port))

			response, err := Marshal(requestMessage,false)

			return response,allocate.PeerAddress,err
		}
	}
	return nil,nil,nil
}

