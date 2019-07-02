package ws

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 30 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type WSConn struct {
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte

	// Buffered channel of inbound messages.
	Recv chan []byte
}

func (c *WSConn) Close() error {
	return c.conn.Close()
}

// Dial: addr should be in the form of host:port
func Dial(addr string) (conn *WSConn, err error) {
	c, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return
	}

	return WrapConn(c), nil
}

// Dial: addr should be in the form of host:port
func WrapConn(c *websocket.Conn) (conn *WSConn) {
	conn = &WSConn{
		conn: c,
		Send: make(chan []byte, 256),
		Recv: make(chan []byte, 256)}

	go conn.writePump()
	go conn.readPump()

	return
}

// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *WSConn) readPump() {
	defer func() {
		c.conn.Close()
	}()
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		mt, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		if mt == websocket.BinaryMessage {
			c.Recv <- message
		}
	}
}

// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *WSConn) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.conn.WriteMessage(websocket.BinaryMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
