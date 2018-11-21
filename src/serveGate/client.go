package serveGate

import (
	"github.com/gorilla/websocket"
	"../pnetWork"

)

type Processor struct{
	client *pnetWork.Client
}

//收到消息
func (p *Processor)OnMessage(conn websocket.Conn,data []byte){

}

//关闭连接
func(p *Processor)Onclose(conn websocket.Conn){

}

//客户连接成功
func (p *Processor)OnClientConnect(conn websocket.Conn){

}

//客户端收到消息
func (p *Processor)OnClientMessage(conn websocket.Conn,data []byte){

}


func Start(){

	client:=pnetWork.Newclient(":12345")
    //注册
	client.Register(&Processor{client:client})
	//开启服务
	client.SocketInit()
}
