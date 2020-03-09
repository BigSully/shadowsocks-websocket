package shadowsocks

import (
	"github.com/georgegloomy/shadowsocks-websocket/websocket"
)

type Conn struct {
	wsConn *ws.WSConn
	*Cipher
}

func NewConn(wsConn *ws.WSConn, cipher *Cipher) (conn *Conn) {
	conn = &Conn{
		wsConn: wsConn,
		Cipher: cipher}

	return
}

func (c *Conn) Close() error {
	return c.wsConn.Close()
}

func (c *Conn) ReadAll() (b []byte, n int, err error) {
	b = <-c.wsConn.Recv // block if no data is available
	n = len(b)

	return
}

func (c *Conn) Write(b []byte) (n int, err error) {
	c.wsConn.Send <- b
	n = len(b)

	return
}
