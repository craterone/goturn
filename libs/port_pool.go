package libs

import (
	"net"
	"strconv"
	"sync"
	"github.com/labstack/gommon/log"
)


type PortPool struct {
	StartPort int
	EndPort int
	Pool map[int]interface{}
	mutex       *sync.Mutex
}

func NewPortPool(start, end int ) *PortPool  {
	log.Printf("start port %d , end	port %d",start,end)
	if (start > 1024) && (end < 65536 ) && (end > start)  {
		pool := make(map[int]interface{})
		p := &PortPool{
			StartPort:start,
			EndPort:end,
			Pool:pool,
			mutex:new(sync.Mutex),
		}
		return p;
	}
	return nil
}

func (p *PortPool)RandSelectPort() int  {
	var port int
	for true {
		port = RandIntRange(p.EndPort,p.StartPort)
		_ ,ok := p.Pool[port]
		if !ok {
			break
		}
	}
	p.SelectPort(port)
	return port
}

func (p *PortPool)SelectPort(port int)  {
	if  port > p.StartPort && port < p.EndPort {

		addr , err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))

		if err != nil {
			Log.Error(err)
		}
		conn , err := net.ListenUDP("udp", addr)
		if err != nil {
			Log.Error(err)
		}else{
			Log.Infof("port %d selected",port)

			p.mutex.Lock()
			p.Pool[port]=true
			p.mutex.Unlock()

			conn.Close()
		}

	}
}

func (p *PortPool)RemovePort(port int)  {
	p.mutex.Lock()
	delete(p.Pool,port)
	p.mutex.Unlock()
}