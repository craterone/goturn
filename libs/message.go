package libs

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type Message struct {
	MessageType   uint16 //2
	MessageLength uint16 //2
	TransID       []byte // 4 + 12
	Attributes    []*Attribute
}


func (m *Message) addAttribute(a *Attribute) {
	m.Attributes = append(m.Attributes, a)
	m.MessageLength += align(a.Length) + 4
}

func (m *Message) addAttributeNoLength(a *Attribute) {
	m.Attributes = append(m.Attributes, a)
}


//UnMarshal creates a Message object from data received by the STUN server
func UnMarshal(data []byte) (*Message, error) {
	length := len(data)
	if length < 20 {
		return nil, ERROR_INVALID_REQUEST
	}

	pkgType := binary.BigEndian.Uint16(data[0:2])

	//check magic cookie
	magicCookieCheck := binary.BigEndian.Uint32(data[4:8]);
	if(magicCookie != magicCookieCheck){
		return nil, ERROR_RFC3489
	}

	msg := new(Message)

	//parse the header
	msg.MessageType = pkgType
	msg.MessageLength = binary.BigEndian.Uint16(data[2:4])

	msg.TransID = data[4:20]

	//if we have leftover data, parse as attributes
	if length > 20 {
		msg.Attributes = make([]*Attribute,0,10)
		i := 20
		for i < length {
			attrType := binary.BigEndian.Uint16(data[i : i+2])
			attrLength := binary.BigEndian.Uint16(data[i+2 : i+4])
			i += 4 + int(attrLength)

			attrValue := data[i-int(attrLength) : i]
			msg.Attributes = append(msg.Attributes,newAttr(attrType,attrValue))

			if pad := int(attrLength) % 4; pad > 0 {
				i += 4 - pad
			}
		}
		//recover here to catch any index errors
		if recover() != nil {
			return nil, ERROR_INVALID_REQUEST
		}
	}
	return msg, nil
}

//Marshal transforms a message into a byte array
func Marshal(m *Message,untilMessageIntegrity bool) ([]byte, error) {
	result := make([]byte, 2048)
	//first do the header
	binary.BigEndian.PutUint16(result[:2], m.MessageType)
	result = append(result[:4], m.TransID...)

	//now we do the attributes
	if m.Attributes != nil {
		i := 20
		for _ , attr := range m.Attributes {
			if untilMessageIntegrity {
				if attr.AttrType == AttributeMessageIntegrity {
					i += 4 + int(attr.Length)
					break
				}
			}

			binary.BigEndian.PutUint16(result[i:i+2], attr.AttrType)
			binary.BigEndian.PutUint16(result[i+2:i+4], attr.Length)

			result = append(result[:i+4], attr.Value...)

			i += 4 + int(attr.Length)
			//if we need to pad, do so
			if pad := int(attr.Length % 4); pad > 0 {
				result = append(result, make([]byte, 4-pad)...)
				i += 4 - pad
			}
		}

		//add length
		binary.BigEndian.PutUint16(result[2:4], uint16(i-20))
	}
	return result, nil
}




func (m Message) hasAttribute(attrType uint16) bool  {
	for _, a := range m.Attributes {
		if a.AttrType == attrType {
			return true
		}
	}
	return false
}

func (m Message) getAttribute(attrType uint16) *Attribute  {
	for _, a := range m.Attributes {
		if a.AttrType == attrType {
			return a
		}
	}
	return nil
}

func (m *Message)setAttribute(attr *Attribute )  {
	for k, v := range m.Attributes {
		if v.AttrType == attr.AttrType {
			m.Attributes[k] = attr
			return
		}
	}
}

func NewResponse(respType uint16,transId []byte,attrs ...*Attribute) *Message {
	respMsg := new(Message)
	respMsg.TransID = transId
	respMsg.MessageType = respType
	respMsg.Attributes = make([]*Attribute,0)

	for _,v := range attrs{
		respMsg.addAttribute(v)
	}

	return respMsg
}



func (m Message) TypeToString() (typeString string)  {
	switch m.MessageType {
	case TypeBindingRequest:
		typeString = "BindRequest"
	case TypeAllocate:
		typeString = "Allocate"
	case TypeBindingResponse:
		typeString = "BindingResponse"
	case TypeAllocateErrorResponse:
		typeString = "AllocateErrorResponse"
	case TypeAllocateResponse:
		typeString = "AllocateResponse"
	case TypeSendIndication:
		typeString = "SendIndication"
	case TypeRefreshRequest:
		typeString = "RefreshRequest"
	default:
		stringByte := make([]byte,2)
		binary.BigEndian.PutUint16(stringByte,m.MessageType)
		typeString = hex.EncodeToString(stringByte)
	}
	return
}

func (m Message) String() string {

	attrString := ""
	if len(m.Attributes) > 0{
		attrString = "\n Attributes : \n"

		for _ , attr := range m.Attributes{
			attrString += attr.String()
		}
	}

	return fmt.Sprintf(`packet : type -> %s , length -> %d , tid -> %X , length of the attr -> %d	%s
			 `,
		m.TypeToString(),m.MessageLength,m.TransID,len(m.Attributes),attrString)
}
