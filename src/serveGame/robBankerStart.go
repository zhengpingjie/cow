package serveGame

import (
	"../common"
	"../serveCenter"
	"fmt"
	"sync"
)

type RobBankerStart struct {
	chairId int32
	robBankerDouble int
	mutex sync.Mutex
}

func RobBankerStartRest(uid string)*RobBankerStart{
	data := common.Manager.RequestData[uid]
	return &RobBankerStart{
		chairId:int32(data.Data["ChairId"].(float64)),
		robBankerDouble:int(data.Data["RobBanker"].(float64)),
	}
}


func (this *RobBankerStart)RobBankerStart(uid string) {
	defer func() {
		this.mutex.Lock()
		common.Wg.Done()
	}()
	tableId := common.UidJoinRoom_.UjR[uid]
	roomType := common.UidJoinRoom_.Urd[uid]
	message:=make(map[string]interface{})
	message["Action"] ="RobBankerStart"
	message["Msg"] = "抢庄成功用户"
	switch roomType {
	case 50:
		room := common.Rooms50_.Room[tableId]
		//设置抢庄状态为true
		common.Rooms50_.Room[tableId] = updateBankerStatus(room,this.chairId,this.robBankerDouble)
		message["UserList"]= getTableUserList(common.Rooms50_.Room[tableId])
		fmt.Printf("\r\nRobBankerStart---------%+v",common.Rooms50_.Room[tableId])
		//向房间每位用户发送抢庄的用户
		serveCenter.SendTableMessage(room,message)
	case 250:
		room := common.Rooms250_.Room[tableId]
		//设置抢庄状态为true
		common.Rooms250_.Room[tableId] = updateBankerStatus(room,this.chairId,this.robBankerDouble)
		message["UserList"]=getTableUserList(common.Rooms50_.Room[tableId])
		//向房间每位用户发送抢庄的用户
		serveCenter.SendTableMessage(room,message)

	case 500:
		room := common.Rooms500_.Room[tableId]
		//设置抢庄状态为true
		common.Rooms500_.Room[tableId] = updateBankerStatus(room,this.chairId,this.robBankerDouble)
		message["UserList"]=getTableUserList(common.Rooms50_.Room[tableId])
		//向房间每位用户发送抢庄的用户
		serveCenter.SendTableMessage(room,message)
	case 1000:
		room := common.Rooms1000_.Room[tableId]
		//设置抢庄状态为true
		common.Rooms1000_.Room[tableId] = updateBankerStatus(room,this.chairId,this.robBankerDouble)
		message["UserList"]=getTableUserList(common.Rooms50_.Room[tableId])
		//向房间每位用户发送抢庄的用户
		serveCenter.SendTableMessage(room,message)
	}
}



//更新抢庄状态
func updateBankerStatus(room common.Room,chairId int32,robBankerDouble int)common.Room {
	for key,user :=range room.UserList{
		if user.ChairId == chairId && user.UserStatus ==1{
			user.Isbanker = true
			user.RobBankerDouble = robBankerDouble
			room.UserList[key] = user
		}
	}
	return room

}


//获取桌面用户列表
func getTableUserList(room common.Room)([]interface{}){
	var UserList []interface{}
	//isBanker := true
	for uid,ok := range room.Uids{
		if ok{
			json_:=make(map[string]interface{})
			json_["UserId"]=common.Manager.UserMapInfo[uid].ID
			json_["UserName"]=common.Manager.UserMapInfo[uid].UserName
			//json_["UserClass"]=common.Manager.UserMapInfo[uid].UserClass
			json_["UserSrc"]=common.Manager.UserMapInfo[uid].UserSrc
			//json_["UserSex"]=common.Manager.UserMapInfo[uid].UserSex
			json_["UserMoney"] = common.Manager.UserMapInfo[uid].UserActivityBalance
			for _,info:=range room.UserList{
				if info.Uid ==uid{
					//if info.Isbanker== false && info.UserStatus ==1{
					//	isBanker = false
					//}
					json_["UserStatus"]=info.UserStatus
					json_["ChairId"]=info.ChairId
					json_["Isbanker"] = info.Isbanker
					json_["RobBankerDouble"] = info.RobBankerDouble
					//break
				}

			}

			UserList=append(UserList, json_)
		}
	}

	return UserList
}
