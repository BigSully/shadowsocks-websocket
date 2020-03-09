package main

import (
	"fmt"
	ss "github.com/georgegloomy/shadowsocks-websocket/shadowsocks"
	"github.com/georgegloomy/shadowsocks-websocket/socks5"
	"github.com/georgegloomy/shadowsocks-websocket/websocket"
	"net"
)

var config map[string]interface{}

var logger ss.Logger

func handleConnection(conn net.Conn) {
	closed := false
	defer func() {
		if !closed {
			conn.Close()
		}
	}()

	var err error = nil
	if err = socks5.HandShake(conn); err != nil { // stage-0
		logger.Println("socks handshake:", err)
		return
	}
	addr, err := socks5.ParseSocksRequest(conn) // stage-1
	rawaddr := addr.Rawaddr
	hostPort := addr.HostPort
	if err != nil {
		logger.Println("error getting request:", err)
		return
	}

	//wsAddr := fmt.Sprintf("ws://%v:%v/", config["server"], config["server_port"])
	wsAddr := config["server"].(string)
	method := config["method"].(string)
	password := config["password"].(string)
	cipher, err := ss.NewCipher(method, password)

	wsConn, err := ws.Dial(wsAddr)
	if err != nil {
		logger.Println("error dialing websocket:", err)
		return
	}

	newConn := ss.NewConn(wsConn, cipher)
	if err != nil {
		logger.Println("error connecting to shadowsocks server:", err)
		return
	}

	// tell the distant server the real server we want to access and help the distant server initialize a new connection
	if _, err = newConn.Write(rawaddr); err != nil {
		newConn.Close()
		return
	}
	logger.Printf("connected to %s via %s from %s\n", hostPort, wsAddr, conn.RemoteAddr().String())

	go ss.PipeNet2WS(conn, *newConn)
	ss.PipeWS2Net(*newConn, conn)
	closed = true
	logger.Println("closed connection to", hostPort)
}

func run(listenAddr string) {
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		logger.Println(err)
	}
	logger.Printf("starting local socks5 server at %v ...\n", listenAddr)
	for {
		conn, err := ln.Accept()
		if err != nil {
			logger.Println("accept:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func main() {
	config = ss.ParseArgs()

	// debug if contains -debug option
	if _, ok := config["debug"]; ok {
		logger = true
	}

	fmt.Printf("%#v\n", config)

	localAddr := fmt.Sprintf("%v:%v", config["local_address"], config["local_port"])

	run(localAddr) // 127.0.0.1:1080
}
