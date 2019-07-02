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
	message := <-c.wsConn.Recv // block if no data is available

	payload := message

	// initialize the decoder with the IV positioned in the beginning of the first packet
	if c.dec == nil {
		iv := make([]byte, c.info.ivLen)
		copy(iv, message[:len(iv)])
		payload = message[len(iv):] // the real payload
		if err = c.initDecrypt(iv); err != nil {
			return
		}
	}

	n = len(payload)
	if n > 0 {
		b = make([]byte, n)
		c.decrypt(b, payload)
	}
	return
}

func (c *Conn) Write(b []byte) (n int, err error) {
	var iv []byte
	if c.enc == nil {
		iv, err = c.initEncrypt()
		if err != nil {
			return
		}
	}

	cipherData := make([]byte, len(b)+len(iv))

	if iv != nil {
		// Put initialization vector in buffer, do a single write to send both
		// iv and data.
		copy(cipherData, iv)
	}

	c.encrypt(cipherData[len(iv):], b)
	c.wsConn.Send <- cipherData
	n = len(cipherData)

	return
}
