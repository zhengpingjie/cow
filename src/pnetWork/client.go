package pnetWork

import (
	"../common"
	"../serveGame"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"sync"
	"time"
)

type Controller struct {
	action  chan common.RequestParam
	register   chan *Client
	unregister chan *Client
}

var ControllerMaster = Controller{
	action:  make(chan common.RequestParam),
	register:   make(chan *Client),
	unregister: make(chan *Client),

}


var mutex sync.Mutex

type Client struct {
	address string
	processor Processor
	delay time.Duration
	socket *websocket.Conn
	send   chan []byte
}



type Message struct {
	Id        string
	Action    string
	Content   string
}

//创建客户端

func Newclient(address string)*Client{
	return &Client{
		address:address,
	}
}

//注册处理
func (C *Client)Register(processor Processor){
	C.processor = processor
}


func(this Controller) start(){
	for{
		select{
			case conn:=<-this.register:
				uid,_:=uuid.NewV4()
				uuidstr := uid.String()
				common.Manager.Clients[uuidstr] = conn.socket
				message:=&Message{Id:uuidstr,Action:"Login",Content:"连接成功"}
				jsonMessage,_:=json.Marshal(message)
				conn.send <- jsonMessage
			case message:=<-this.action:
				 mainController(message)

		}
	}

}

//读取消息
func (C *Client)read(){
	for{
		_,message,err:=C.socket.ReadMessage();

		if err !=nil{
			goto ERR
			break
		}
		var param = &common.RequestParam{}
		json.Unmarshal(message,param)

		//mainController(*param)
		ControllerMaster.action <- *param
	}

	ERR:
		C.socket.Close()

}




//写消息
func(C *Client)write(){
	defer func() {
		C.socket.Close()
	}()

	for {
		select {
		case message, ok := <-C.send:
			if !ok {
				C.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			C.socket.WriteMessage(websocket.TextMessage, message)
		}
	}

}




func (C *Client)SocketInit(){
	log.Print("websocket连接服务启动")
	go ControllerMaster.start()
	http.HandleFunc("/ws",tosocket)
	err:=http.ListenAndServe(C.address,nil)
	if err != nil {
		log.Fatal("ListenAndServe", err.Error())
	}

}

func tosocket(w http.ResponseWriter,r *http.Request){
	conn,err :=(&websocket.Upgrader{
		// 读取存储空间大小
		ReadBufferSize:1024,
		// 写入存储空间大小
		WriteBufferSize:1024,
		CheckOrigin: func(r *http.Request) bool {return true} }).Upgrade(w,r,nil)
		if err !=nil{
			log.Print(err)
			return
		}

	client := &Client{ socket: conn, send: make(chan []byte)}
	ControllerMaster.register <- client
	go client.read()
	go client.write()
}



//action总调度器
func  mainController(data common.RequestParam){
	//异常处理
	defer func() {
		if err:=recover();err !=nil{
			fmt.Println(err)
		}
	}()

	//【Lock方法将rw锁定为写入状态，禁止其他线程读取或者写入】
	getValue:=reflect.ValueOf(ControllerMaster)
	if data.Uid ==""{
		return;
	}
	mutex.Lock()
	common.Manager.RequestData[data.Uid] = data
	mutex.Unlock()
	methodValue := getValue.MethodByName(data.Action)
	args := []reflect.Value{reflect.ValueOf(data.Uid)}
	methodValue.Call(args)
}

//方法名MethodByName调用方法必须是大写的，否则会抛异常
func (this Controller)Register(uid string){
	common.Wg.Add(1)
	contro:=serveGame.Rest(uid)
	go contro.UserRegister(uid)
	common.Wg.Wait()
	println("Register注册--------------")
	runtime.Gosched()
}


//登陆
func (this Controller)Login(uid string){
	common.Wg.Add(1)
	contro := serveGame.LoginRest(uid)
	go contro.UserLogin(uid)
	common.Wg.Wait()
	println("Login登陆--------------")
	runtime.Gosched()
}


//选择场次
func (this Controller)SelectRound(uid string){
	//验证token
	mutex.Lock()
	isok:=common.IsToken(uid,common.Manager.RequestData[uid].Token)
	mutex.Unlock()

	if !isok{
		return
	}
	contro := serveGame.SelectRoundRest(uid)
	common.Wg.Add(1)
	go contro.UserSelectRound(uid)
	common.Wg.Wait()
	println("SelectRound选择场次--------------")
	runtime.Gosched()
}

//抢庄开始准备 uid chairId //
func (this Controller)RobBankerReady(uid string){
	mutex.Lock()
	defer mutex.Unlock()
	contro :=serveGame.RobBankerReadyRest(uid)
	contro.RobBankerReady(uid)
}

////更新抢庄倒计时【时间】
func(this Controller)UpdateRoomCountDownTime(uid string){
	common.Wg.Add(1)
	contro :=serveGame.CountDownTimeRest(uid)
	go contro.UpdateCountDownTime(uid)
	common.Wg.Wait()
	println("UpdateRoomCountDownTime更新抢庄倒计时--------------")
	runtime.Gosched()
}

////抢庄倒计时结束时
func(this Controller)RoomCountDownTimeOver(uid string){
	common.Wg.Add(1)
	contro :=serveGame.TimeOverRest()
	go contro.TimeOver(uid)
	common.Wg.Wait()
	println("RoomCountDownTimeOver抢庄倒计时结束时--------------")
	runtime.Gosched()
}

//开始抢庄uid chairId

func (this Controller)RobBankerStart(uid string){
	common.Wg.Add(1)
	contro :=serveGame.RobBankerStartRest(uid)
	go contro.RobBankerStart(uid)
	common.Wg.Wait()
	println("RobBankerStart开始抢庄--------------")
}
//抢庄结束
func(this Controller)RobBankerOver(uid string){
	common.Wg.Add(1)
	contro :=serveGame.RobBankerOverRest()
	go contro.RobBankerOver(uid)
	common.Wg.Wait()
	println("RobBankerOver抢庄结束--------------")
	runtime.Gosched()
}

//下注开始
func(this Controller)BetStart(uid string){
	common.Wg.Add(1)
	contro :=serveGame.BetStartRest(uid)
	go contro.BetStart(uid)
	common.Wg.Wait()
	println("BetStart下注开始--------------")
	runtime.Gosched()
}

//下注结束
func(this Controller)BetOver(uid string){
	common.Wg.Add(1)
	contro :=serveGame.BetOverRest()
	go contro.BetOver(uid)
	common.Wg.Wait()
	println("BetOver下注结束--------------")
	runtime.Gosched()
}


//发牌开始
func(this Controller)SendPokerStart(uid string){
	common.Wg.Add(1)
	contro :=serveGame.SendPokerStartRest(uid)
	go contro.SendPokerStart(uid)
	common.Wg.Wait()
	println("SendPokerStart发牌开始--------------")
	runtime.Gosched()
}

//摊牌
func(this Controller)ShowPoker(uid string){
	common.Wg.Add(1)
	contro :=serveGame.ShowPokerRest(uid)
	go contro.ShowPoker(uid)
	common.Wg.Wait()
	println("SendPokerStart发牌开始--------------")
	runtime.Gosched()
}

