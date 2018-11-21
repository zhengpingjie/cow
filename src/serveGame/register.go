package serveGame

import (
	"../common"
	"encoding/json"
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

type Register struct {
	username string
	userPhone string
	password string
	mutex sync.Mutex
}

//var mutex sync.Mutex
//注册struct
func Rest(uid string)*Register{
	data := common.Manager.RequestData[uid]
	return &Register{
		username:data.Data["userName"].(string),
		userPhone:data.Data["userPhone"].(string),
		password:data.Data["password"].(string),
	}
}
//用户注册 返回值success 【0注册失败，1注册成功，2已经注册】
func (this *Register)UserRegister(uid string){

	this.mutex.Lock()
	defer func() {
		this.mutex.Unlock()
		common.Wg.Done()
	}()

	if err :=DB.Ping();err !=nil{
		panic(err)
		return
	}
	row:=DB.QueryRow("select UserName from pre_user where UserPhone = ?",this.userPhone)
	var UserName string
	row.Scan(&UserName)
	var jsonMessage =make(map[string]string)
	if UserName==""{
		tx,_:= DB.Begin()
		result,_:=tx.Exec("INSERT INTO pre_user(UserName,UserPhone,UserPassword,AddTime)values(?,?,?,?)",this.username,this.userPhone,this.password,time.Unix(time.Now().Unix(),0))
		err:=tx.Commit()

		if err !=nil{
			tx.Rollback()
		}

		if ok,_:=result.LastInsertId();ok>0{
			//fmt.Println(ok)
			jsonMessage["Success"]="1"
			jsonMessage["Msg"]="注册成功"
			jsonMessage["Action"] ="Login"
		}else{
			jsonMessage["Success"]="0"
			jsonMessage["Msg"]="注册失败"
			jsonMessage["Action"] ="Register"
		}
	}else{
		jsonMessage["Success"]="2"
		jsonMessage["Msg"]="已经注册，请登录"
		jsonMessage["Action"] ="Login"
	}
	message,_:=json.Marshal(jsonMessage)
	common.Manager.Clients[uid].WriteMessage(websocket.TextMessage,message)
}

