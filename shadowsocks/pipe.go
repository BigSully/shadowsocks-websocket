package shadowsocks

import (
	"log"
	"net"
)

func PipeNet2WS(src net.Conn, dst Conn) {
	defer dst.Close()
	buf := leakyBuf.Get()
	defer leakyBuf.Put(buf)
	for {
		n, err := src.Read(buf)
		if n > 0 {
			if _, err := dst.Write(buf[0:n]); err != nil {
				log.Println("write:", err)
				break
			}
		}
		if err != nil {
			break
		}
	}
	return
}

func PipeWS2Net(src Conn, dst net.Conn) {
	defer dst.Close()
	for {
		buf, n, err := src.ReadAll()
		if n > 0 {
			if _, err := dst.Write(buf); err != nil {
				log.Println("write:", err)
				break
			}
		}
		if err != nil {
			break
		}
	}
	return
}
