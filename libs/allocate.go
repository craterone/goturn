package libs

import (
	"net"
	"time"
)




type Allocate struct {
	Token string
	ClientAddress *net.UDPAddr
	PeerAddress *net.UDPAddr
	Relay *Relay

	ExpiresTime int
	ConnectTime uint32
	IsExpired bool
	ExpiresTicker *time.Ticker `json:"-"`

	BytesRecv int
	BytesSend int
}

func (a *Allocate)TimerRun()  {
	go func() {
		for range a.ExpiresTicker.C {
			if a.ExpiresTime <= 0 {
				a.IsExpired = true
				a.ExpiresTicker.Stop()

				AllocateMutex.Lock()
				delete(GlobalAllocates,a.ClientAddress.String())
				AllocateMutex.Unlock()
				break
			}else{
				a.ConnectTime++;
				a.ExpiresTime--;
			}
		}
	}()
}

//func AllocatesMap2JSON(allocates map[string]*Allocate) []byte {
//
//}