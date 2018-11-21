package serveGame

import (
	"../common"
	"encoding/json"
	"github.com/gorilla/websocket"
	"sync"
)

type Login struct {
	userPhone string
	password string
	sign string  //签名
	datetime string //时间戳
	userName string
	mutex sync.Mutex
}
//注册struct
func LoginRest(uid string)*Login{
	data := common.Manager.RequestData[uid]
	return &Login{
		userPhone:data.Data["userPhone"].(string),
		password:data.Data["password"].(string),
		sign:data.Data["sign"].(string),
		datetime:data.Data["datetime"].(string),
	}
}

//用户注册 返回值success 【0登陆失败，1登陆成功】
func (this *Login)UserLogin(uid string){
	this.mutex.Lock()
	defer func() {
		this.mutex.Unlock()
		common.Wg.Done()
	}()
	row:=DB.QueryRow("select UserActivityBalance,UserSafeBalance,ID,UserName,UserPhone,UserBirthday,UserClass,UserSex from pre_user where UserPhone = ? and UserPassword = ?",this.userPhone,this.password)
	var ID,UserClass,UserSex int
	var UserName,UserPhone,UserBirthday string
	var UserSafeBalance float64 //用户保险柜余额
	var UserActivityBalance float64 //用户活动余额
	row.Scan(&UserActivityBalance,&UserSafeBalance,&ID,&UserName,&UserPhone,&UserBirthday,&UserClass,&UserSex)
	//row.Scan(&UserSafeBalance,&UserActivityBalance)
	if ID < 1{
		common.Response.Action = "Login"
		common.Response.Uid = uid
		common.Response.Success = "0"
		common.Response.Msg = "用户名或者密码错误"
		message,_:=json.Marshal(common.Response)
		common.Manager.Clients[uid].WriteMessage(websocket.TextMessage,message)
		return
	}else{
		common.Response.Action = "SelectRound"
		common.Response.Uid = uid
		common.Response.Success = "1"
		common.Response.Msg = "登陆成功"
		//获取MD5加密字符串
		common.Response.Token =common.GetMd5(this.datetime,this.sign)

		common.UserInfo.ID = ID
		common.UserInfo.UserName = UserName
		common.UserInfo.UserPhone = UserPhone
		common.UserInfo.UserBirthday = UserBirthday
		common.UserInfo.UserClass = UserClass
		common.UserInfo.UserSex = UserSex
		common.UserInfo.UserActivityBalance = UserActivityBalance
		common.UserInfo.UserSafeBalance = UserSafeBalance



		//登陆成功 添加用户到游戏中心
		common.Manager.UserMapInfo[uid] = common.UserInfo
		//获取当前在线总人数
		common.Response.OnlineNum = len(common.Manager.UserMapInfo)
		common.Response.Data = common.UserInfo

		//存储用户签名时间戳
		common.Token.DateTime[uid] = this.datetime
		common.Token.Sign[uid] = this.sign

	}
	message,_:=json.Marshal(common.Response)
	//fmt.Printf("%v",common.Manager.UserMapInfo)
	common.Manager.Clients[uid].WriteMessage(websocket.TextMessage,message)
}


