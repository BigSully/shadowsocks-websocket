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
				log.Println("PipeNet2WS write:", err)
				break
			}
		}
		if err != nil {
			log.Println("PipeNet2WS Read:", err)
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
				log.Println("PipeWS2Net write:", err)
				break
			}
		}
		if err != nil {
			log.Println("PipeWS2Net ReadAll:", err)
			break
		}
	}
	return
}
