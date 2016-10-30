package libs


import (
	"log"
	"crypto/rand"
	mrand "math/rand"
	"net"
	"errors"
	"strings"
	"strconv"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/binary"
	"time"
	"unsafe"
)

func HmacSha1(value,key []byte) []byte {
	hasher := hmac.New(sha1.New,key)
	hasher.Write(value)
	digest := hasher.Sum(nil)
	return digest
}


func RandBytes(length int) (r []byte) {
	if length < 64 {
		r = make([]byte, length)
		_, err := rand.Read(r)

		if err != nil {
			log.Panicln(err)
		}
	}else {
		log.Panicf("the max length of randbyte is 64 , %d not supported \n",length)
	}
	return
}

func PrintModuleLoaded(moduleName string)  {
	log.Printf("< %s > module loads successfully",moduleName)
}

func PrintModuleRelease(moduleName string)  {
	log.Printf("< %s > module releases successfully",moduleName)
}

func HostIP() (ipAddress string, err error) {
	var ifaces []net.Interface
	ifaces, err = net.Interfaces()
	if err != nil {
		return
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		var addrs []net.Addr
		addrs, err = iface.Addrs()
		if err != nil {
			return
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			ipAddress = ip.String()
			return
		}
	}
	err = errors.New("are you connected to the network?")
	return
}

func IsValidIPv4(host string) bool {
	parts := strings.Split(host, ".")

	if len(parts) < 4 {
		return false
	}

	for _,x := range parts {
		if i, err := strconv.Atoi(x); err == nil {
			if i < 0 || i > 255 {
				return false
			}
		} else {
			return false
		}

	}
	return true
}

func RandIntRange(max,min int) int   {
	mrand.Seed(time.Now().UnixNano())
	num := mrand.Intn(max - min) + min
	Log.Infof("random :%d",num)
	return num
}

func parseAddressV4(strAddress string) *net.UDPAddr {
	arrAddr := strings.Split(strAddress,":")
	port , _ := strconv.Atoi(arrAddr[1])
	address := new(net.UDPAddr)
	address.IP = net.ParseIP(arrAddr[0])
	address.Port = port
	return address
}

func generateTransactionID() []byte  {
	transID := make([]byte, 16)
	binary.BigEndian.PutUint32(transID[:4], magicCookie)
	rand.Read(transID[4:])
	return transID
}

func str2bytes(s string)[] byte  {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0],x[1],x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}