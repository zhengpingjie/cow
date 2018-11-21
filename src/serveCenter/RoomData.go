package serveCenter

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"../common"
)


//创建50桌子
func Createtable50(num int)string {
	Room_ := common.Room{Uids: make(map[string]bool)}
	tableNum := num + 1
	TableId := fmt.Sprintf("%s%d", "00", tableNum)
	Room_.TableStatus = 0
	Room_.TableId = TableId
	Room_.UserNum = 1
	Room_.UserActivity = 1
	common.Rooms50_.Room[TableId] = Room_
	return TableId
}

//添加用户到50桌子
func AddUserToTable50(uid string){
	Room_ := common.Room{Uids:make(map[string]bool)}
	GameUser_ := common.GameUser{}
	//椅子编号为房间人数
	tableId := common.UidJoinRoom_.UjR[uid]
	GameUser_.ChairId = common.Rooms50_.Room[tableId].UserNum
	GameUser_.Uid = uid
	GameUser_.UserStatus = 1
	//设置默认值
	GameUser_.RobBankerDouble = 1
	GameUser_.Bet = 5

	Room_.UserList=append(Room_.UserList, GameUser_)
	Room_.Uids[uid] = true
	Room_.TableId = common.Rooms50_.Room[tableId].TableId
	Room_.TableStatus = common.Rooms50_.Room[tableId].TableStatus
	Room_.UserNum = common.Rooms50_.Room[tableId].UserNum
	Room_.TableId = common.Rooms50_.Room[tableId].TableId
	Room_.Playing = common.Rooms50_.Room[tableId].Playing
	Room_.UserActivity = common.Rooms50_.Room[tableId].UserActivity
	//添加用户到房间
	common.Rooms50_.Room[tableId]= Room_
}


//创建250桌子
func Createtable250(num int)string{
	Room_ := common.Room{Uids:make(map[string]bool)}
	tableNum:=num+1
	TableId := fmt.Sprintf("%s%d","00",tableNum)
	Room_.TableStatus = 0
	Room_.TableId = TableId
	Room_.UserNum = 1
	common.Rooms250_.Room[TableId] = Room_
	return TableId
}

//添加用户到250桌子
func AddUserToTable250(uid string){
	Room_ := common.Room{Uids:make(map[string]bool)}
	GameUser_ := common.GameUser{}
	//椅子编号为房间人数
	tableId := common.UidJoinRoom_.UjR[uid]
	GameUser_.ChairId = common.Rooms50_.Room[tableId].UserNum
	GameUser_.Uid = uid
	Room_.UserList=append(Room_.UserList, GameUser_)
	Room_.Uids[uid] = true
	//添加用户到房间
	common.Rooms250_.Room[tableId]= Room_
}

//创建500桌子
func Createtable500(num int)string{
	Room_ := common.Room{Uids:make(map[string]bool)}
	tableNum:=num+1
	TableId := fmt.Sprintf("%s%d","00",tableNum)
	Room_.TableStatus = 0
	Room_.TableId = TableId
	Room_.UserNum = 1
	common.Rooms500_.Room[TableId] = Room_
	return TableId
}

//添加用户到500桌子
func AddUserToTable500(uid string){
	Room_ := common.Room{Uids:make(map[string]bool)}
	GameUser_ := common.GameUser{}
	//椅子编号为房间人数
	tableId := common.UidJoinRoom_.UjR[uid]
	GameUser_.ChairId = common.Rooms50_.Room[tableId].UserNum
	GameUser_.Uid = uid
	Room_.UserList=append(Room_.UserList, GameUser_)
	Room_.Uids[uid] = true
	//添加用户到房间
	common.Rooms250_.Room[tableId]= Room_
}


//创建1000桌子
func Createtable1000(num int)string{
	Room_ := common.Room{Uids:make(map[string]bool)}
	tableNum:=num+1
	TableId := fmt.Sprintf("%s%d","00",tableNum)
	Room_.TableStatus = 0
	Room_.TableId = TableId
	Room_.UserNum = 1
	common.Rooms1000_.Room[TableId] = Room_
	return TableId
}

//添加用户到1000桌子
func AddUserToTable1000(uid string){
	Room_ := common.Room{Uids:make(map[string]bool)}
	GameUser_ := common.GameUser{}
	//椅子编号为房间人数
	tableId := common.UidJoinRoom_.UjR[uid]
	GameUser_.ChairId = common.Rooms50_.Room[tableId].UserNum
	GameUser_.Uid = uid

	Room_.UserList=append(Room_.UserList, GameUser_)
	Room_.Uids[uid] = true
	//添加用户到房间
	common.Rooms250_.Room[tableId]= Room_
}

//发送消息给房间每个用户
func SendTableMessage(room common.Room,RoomInfo map[string]interface{}){
	//获取用户列表
	for uid,ok:=range room.Uids{
		if ok {
			jsonMessage,_:=json.Marshal(RoomInfo)
			common.Manager.Clients[uid].WriteMessage(websocket.TextMessage,jsonMessage)
		}

	}

}

//发送消息给单个用户
func SendClientMessage(uid string,message map[string]interface{}){
	jsonMessage,_:= json.Marshal(message)
	common.Manager.Clients[uid].WriteMessage(websocket.TextMessage,jsonMessage)
}

//获取桌面用户列表
func GetTableUserList(room common.Room)[]interface{}{
	fmt.Printf("roomCenter148-------------%+v\r\n",room)
	fmt.Printf("roomCenter149-------------%+v\r\n",common.Manager.UserMapInfo)
	var UserList []interface{}
	for uid,ok := range room.Uids{
		if ok{
			json_:=make(map[string]interface{})
			json_["UserId"]=common.Manager.UserMapInfo[uid].ID
			json_["UserName"]=common.Manager.UserMapInfo[uid].UserName
			//json_["UserClass"]=common.Manager.UserMapInfo[uid].UserClass
			json_["UserSrc"]=common.Manager.UserMapInfo[uid].UserSrc
			json_["UserSex"]=common.Manager.UserMapInfo[uid].UserSex
			json_["UserMoney"] = common.Manager.UserMapInfo[uid].UserActivityBalance
			for _,info:=range room.UserList{
				fmt.Println("------roomCenter155---",info.ChairId)
				if info.Uid ==uid{
					json_["UserStatus"]=info.UserStatus
					json_["ChairId"]=info.ChairId
					//break
				}

			}

			UserList=append(UserList, json_)
		}
	}

	return UserList
}


//清除房间用户
func ClearRoomClients(uid string,room common.Room)common.Room{
	for key,userinfo := range room.UserList{
		if userinfo.Uid == uid{
			//移除房间用户
			room.UserList = append(room.UserList[key:],room.UserList[key+1:]...)
		}
	}
	room.UserNum = room.UserNum-1
	return room
}