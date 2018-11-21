package serveGame

import (
	"../common"
	"../serveCenter"
	"sync"
)

type SendPokerStart struct {
	chairId int32
	mutex sync.Mutex
}

func SendPokerStartRest(uid string)*SendPokerStart{
	data := common.Manager.RequestData[uid]
	return &SendPokerStart{
		chairId:int32(data.Data["ChairId"].(float64)),
	}
}


func (this *SendPokerStart)SendPokerStart(uid string){
	this.mutex.Lock()
	defer func() {
		this.mutex.Unlock()
		common.Wg.Done()
	}()
	tableId := common.UidJoinRoom_.UjR[uid]
	roomType := common.UidJoinRoom_.Urd[uid]

	switch roomType {
	case 50:
		room := common.Rooms50_.Room[tableId]
		//进入摊牌阶段
		room.Playing = 4
		if len(room.CartList) <=0 {
			room.CartList = common.GetRandNumber(1,52,int(room.UserActivity*5))
		}
		common.Rooms50_.Room[tableId] = room
		//获取用户牌
		common.Rooms50_.Room[tableId]=getUserPoker(this.chairId,room)
		//发牌
		message := getMessage(this.chairId,common.Rooms50_.Room[tableId])
		message["TimeOut"] = common.ShowPokerCountDown
		serveCenter.SendClientMessage(uid,message)
	case 250:
		room := common.Rooms250_.Room[tableId]
		//进入摊牌阶段
		room.Playing = 4
		if len(room.CartList) <=0 {
			room.CartList = common.GetRandNumber(1,52,10)
		}
		common.Rooms250_.Room[tableId] = room
		//获取用户牌
		common.Rooms50_.Room[tableId]=getUserPoker(this.chairId,room)
		//发牌
		message := getMessage(this.chairId,common.Rooms50_.Room[tableId])
		message["TimeOut"] = common.ShowPokerCountDown
		serveCenter.SendClientMessage(uid,message)
	case 500:
		room := common.Rooms500_.Room[tableId]
		//进入摊牌阶段
		room.Playing = 4
		if len(room.CartList) <=0 {
			room.CartList = common.GetRandNumber(1,52,10)
		}
		common.Rooms500_.Room[tableId] = room
		//获取用户牌
		common.Rooms50_.Room[tableId]=getUserPoker(this.chairId,room)
		//发牌
		message := getMessage(this.chairId,common.Rooms50_.Room[tableId])
		message["TimeOut"] = common.ShowPokerCountDown
		serveCenter.SendClientMessage(uid,message)
	case 1000:
		room := common.Rooms1000_.Room[tableId]
		//进入摊牌阶段
		room.Playing = 4
		if len(room.CartList) <=0 {
			room.CartList = common.GetRandNumber(1,52,10)
		}
		common.Rooms1000_.Room[tableId] = room
		//获取用户牌
		common.Rooms1000_.Room[tableId]=getUserPoker(this.chairId,room)
		//发牌
		message := getMessage(this.chairId,common.Rooms1000_.Room[tableId])
		message["TimeOut"] = common.ShowPokerCountDown
		serveCenter.SendClientMessage(uid,message)
	}
}

//获取用户牌
func getUserPoker(chairId int32,room common.Room)common.Room{
	start:=chairId*5-5
	end := chairId*5
	pokers := room.CartList[start:end]
	for key,user := range room.UserList{
		if user.ChairId == chairId{
			user.UserCartList = pokers
			room.UserList[key] = user
		}
	}

	return room
}

//获取消息
func getMessage(chairId int32,room common.Room)map[string]interface{}{
	message:=make(map[string]interface{})
	UserList:=make(map[string]interface{})
	var UserLists []map[string]interface{}
	message["Action"] ="ShowPoker"
	message["Msg"] = "进入摊牌阶段"
	for _,user := range room.UserList{
		UserList["UserId"] = common.Manager.UserMapInfo[user.Uid].ID
		UserList["UserName"]=common.Manager.UserMapInfo[user.Uid].UserName
		UserList["UserSrc"]=common.Manager.UserMapInfo[user.Uid].UserSrc
		UserList["UserMoney"]= common.Manager.UserMapInfo[user.Uid].UserActivityBalance
		if user.ChairId == chairId{
			UserList["CardList"] = user.UserCartList
		}else{
			UserList["CardList"] = []int{0,0,0,0,0}
		}
		UserLists = append(UserLists,UserList)
	}
	message["UserList"] = UserLists
	return message
}


