package main

import (
	"encoding/binary"
	"fmt"
	ss "github.com/georgegloomy/shadowsocks-websocket/shadowsocks"
	"github.com/georgegloomy/shadowsocks-websocket/websocket"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"syscall"
)

const (
	idType  = 0 // address type index
	idIP0   = 1 // ip address start index
	idDmLen = 1 // domain address length index
	idDm0   = 2 // domain address start index

	typeIPv4 = 1 // type is ipv4 address
	typeDm   = 3 // type is domain address
	typeIPv6 = 4 // type is ipv6 address

	lenIPv4   = net.IPv4len + 2 // ipv4 + 2port
	lenIPv6   = net.IPv6len + 2 // ipv6 + 2port
	lenDmBase = 2               // 1addrLen + 2port, plus addrLen
	// lenHmacSha1 = 10

	AddrMask byte = 0xf
)

var config map[string]interface{}
var logger ss.Logger

func sanitizeAddr(addr net.Addr) string {
	return addr.String()
}

var cipher *ss.Cipher

func handleConnection(wconn *websocket.Conn) {

	logger.Printf("new client %s->%s\n", sanitizeAddr(wconn.RemoteAddr()), wconn.LocalAddr())
	newConn := ws.WrapConn(wconn)
	conn := ss.NewConn(newConn, cipher.Copy())

	var host string

	closed := false
	defer func() {
		log.Printf("closed pipe %s<->%s\n", sanitizeAddr(wconn.RemoteAddr()), host)
		if !closed {
			conn.Close()
		}
	}()

	host, err := extractUpstreamServer(conn) // stage-0
	if err != nil {
		log.Println("error getting request", sanitizeAddr(wconn.RemoteAddr()), wconn.LocalAddr(), err)
		closed = true
		return
	}
	// ensure the host does not contain some illegal characters, NUL may panic on Win32
	if strings.ContainsRune(host, 0x00) {
		log.Println("invalid domain name.")
		closed = true
		return
	}

	logger.Println("connecting", host)
	remote, err := net.Dial("tcp", host)
	if err != nil {
		if ne, ok := err.(*net.OpError); ok && (ne.Err == syscall.EMFILE || ne.Err == syscall.ENFILE) {
			// log too many open file error
			// EMFILE is process reaches open file limits, ENFILE is system limit
			log.Println("dial error:", err)
		} else {
			log.Println("error connecting to:", host, err)
		}
		return
	}
	defer func() {
		if !closed {
			remote.Close()
		}
	}()

	//logger.Printf("piping %s<->%s", sanitizeAddr(conn.RemoteAddr()), host)

	go ss.PipeWS2Net(*conn, remote) // stage-5
	ss.PipeNet2WS(remote, *conn)

	closed = true
	return
}

func extractUpstreamServer(conn *ss.Conn) (host string, err error) {
	//ss.SetReadTimeout(conn)

	// buf size should at least have the same size with the largest possible
	// request size (when addrType is 3, domain name has at most 256 bytes)
	// 1(addrType) + 1(lenByte) + 255(max length address) + 2(port) + 10(hmac-sha1)
	//buf := make([]byte, 269)
	//// read till we get possible domain length field
	//if _, err = io.ReadFull(conn, buf[:idType+1]); err != nil {
	//	return
	//}

	buf, n, error := conn.ReadAll()
	if error != nil {
		return
	}
	if n < idType+1 {
		return
	}
	var reqStart, reqEnd int
	addrType := buf[idType]
	switch addrType & AddrMask {
	case typeIPv4:
		reqStart, reqEnd = idIP0, idIP0+lenIPv4
	case typeIPv6:
		reqStart, reqEnd = idIP0, idIP0+lenIPv6
	case typeDm:
		if n < idDmLen+1 {
			return
		}
		_, reqEnd = idDm0, idDm0+int(buf[idDmLen])+lenDmBase
	default:
		err = fmt.Errorf("addr type %d not supported", addrType&AddrMask)
		return
	}

	if n < reqStart || n < reqEnd {
		return
	}

	// Return string for typeIP is not most efficient, but browsers (Chrome,
	// Safari, Firefox) all seems using typeDm exclusively. So this is not a
	// big problem.
	switch addrType & AddrMask {
	case typeIPv4:
		host = net.IP(buf[idIP0 : idIP0+net.IPv4len]).String()
	case typeIPv6:
		host = net.IP(buf[idIP0 : idIP0+net.IPv6len]).String()
	case typeDm:
		host = string(buf[idDm0 : idDm0+int(buf[idDmLen])])
	}
	// parse port
	port := binary.BigEndian.Uint16(buf[reqEnd-2 : reqEnd])
	host = net.JoinHostPort(host, strconv.Itoa(int(port)))
	return
}

var upgrader = websocket.Upgrader{} // use default options
func serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	handleConnection(conn)
}

func main() {
	config = ss.ParseArgs()

	// debug if contains -debug option
	if _, ok := config["debug"]; ok {
		logger = true
	}

	fmt.Printf("%#v\n", config)

	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello!"))
	})
	http.HandleFunc("/", serveWs)

	method := config["method"].(string)
	password := config["password"].(string)
	cipher, _ = ss.NewCipher(method, password)
	serverPort := fmt.Sprintf(":%v", config["server_port"])
	if err := http.ListenAndServe(serverPort, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
