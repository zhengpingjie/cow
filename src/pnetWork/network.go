package pnetWork

import "github.com/gorilla/websocket"

type Processor interface {
	OnMessage(conn websocket.Conn,data []byte)
	Onclose(conn websocket.Conn)
	OnClientConnect(conn websocket.Conn)
	OnClientMessage(conn websocket.Conn,data []byte)
}