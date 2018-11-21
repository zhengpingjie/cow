package serveGame

import (
	"../common"
	"../serveCenter"
	"sync"
)

type BetStart struct {
	chairId int32
	BetStartDouble int
	mutex sync.Mutex

}

func BetStartRest(uid string)*BetStart{
	data := common.Manager.RequestData[uid]
	return &BetStart{
		chairId:int32(data.Data["ChairId"].(float64)),
		BetStartDouble:int(data.Data["BetStartDouble"].(float64)),
	}
}

func(this *BetStart)BetStart(uid string){
	this.mutex.Lock()
	defer func() {
		this.mutex.Unlock()
		common.Wg.Done()
	}()
	tableId := common.UidJoinRoom_.UjR[uid]
	roomType := common.UidJoinRoom_.Urd[uid]
	message:=make(map[string]interface{})
	message["Action"] ="BetStart"
	message["Msg"] = "下注成功用户"
	switch roomType {
	case 50:
		room := common.Rooms50_.Room[tableId]
		common.Rooms50_.Room[tableId] = updateBetStatus(room,this.chairId,this.BetStartDouble)
		message["BankerUserId"] = common.Manager.UserMapInfo[room.BankerUid].ID
		message["UserList"] = getTableBetUserList(common.Rooms50_.Room[tableId])
		//向房间每位用户发送下注的用户
		serveCenter.SendTableMessage(room,message)
	case 250:
		room := common.Rooms250_.Room[tableId]
		common.Rooms250_.Room[tableId] = updateBetStatus(room,this.chairId,this.BetStartDouble)
		message["BankerUserId"] = common.Manager.UserMapInfo[room.BankerUid].ID
		message["UserList"] = getTableBetUserList(common.Rooms250_.Room[tableId])
		//向房间每位用户发送下注的用户
		serveCenter.SendTableMessage(room,message)
	case 500:
		room := common.Rooms500_.Room[tableId]
		common.Rooms500_.Room[tableId] = updateBetStatus(room,this.chairId,this.BetStartDouble)
		message["BankerUserId"] = common.Manager.UserMapInfo[room.BankerUid].ID
		message["UserList"] = getTableBetUserList(common.Rooms500_.Room[tableId])
		//向房间每位用户发送下注的用户
		serveCenter.SendTableMessage(room,message)
	case 1000:
		room := common.Rooms1000_.Room[tableId]
		common.Rooms1000_.Room[tableId] = updateBetStatus(room,this.chairId,this.BetStartDouble)
		message["BankerUserId"] = common.Manager.UserMapInfo[room.BankerUid].ID
		message["UserList"] = getTableBetUserList(common.Rooms1000_.Room[tableId])
		//向房间每位用户发送下注的用户
		serveCenter.SendTableMessage(room,message)
	}
}


//更新下注状态
func updateBetStatus(room common.Room,chairId int32,BetStartDouble int)common.Room {
	for key,user :=range room.UserList{
		if user.ChairId == chairId && user.UserStatus ==1{
			user.Bet = BetStartDouble
			user.IsBet = true
			room.UserList[key] = user
		}
	}
	//游戏进入押注阶段
	return room

}




