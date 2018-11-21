package serveGame

import (
	"../common"
	"../serveCenter"
	"sync"
)

type RobBankerOver struct {
	mutex sync.Mutex
}

func RobBankerOverRest()*RobBankerOver{
	return &RobBankerOver{

	}
}

func(this *RobBankerOver)RobBankerOver(uid string){
	this.mutex.Lock()
	defer func() {
		this.mutex.Unlock()
		common.Wg.Done()
	}()

	tableId := common.UidJoinRoom_.UjR[uid]
	roomType := common.UidJoinRoom_.Urd[uid]
	message:=make(map[string]interface{})
	message["Action"] ="BetStart"
	message["Msg"] ="下注开始"
	message["TimeOut"] = common.BetCountDown
	switch roomType {
	case 50:
		room := common.Rooms50_.Room[tableId]
		if room.TableStatus==2 && room.Playing == 2{
			//已经有用户发送了倒计时结束请求return
			return
		}
		//设置抢庄状态为true
		Banker :=getBankerPerson(room)
		room.BankerUid = Banker["Uid"].(string)
		room.Playing = 2
		common.Rooms50_.Room[tableId] = room
		message["BankerUserId"] = Banker["UserId"]
		message["UserList"]=getTableUserList(common.Rooms50_.Room[tableId])
		message["Banker"] = Banker
		//向房间每位用户发送抢庄的用户
		serveCenter.SendTableMessage(room,message)
	case 250:
		room := common.Rooms250_.Room[tableId]
		if room.TableStatus==2 && room.Playing == 2{
			//已经有用户发送了倒计时结束请求return
			return
		}
		//设置抢庄状态为true
		Banker :=getBankerPerson(room)
		room.BankerUid = Banker["Uid"].(string)
		room.Playing = 2
		common.Rooms250_.Room[tableId] = room
		message["UserList"]=getTableUserList(common.Rooms250_.Room[tableId])
		message["Banker"] = Banker
		//向房间每位用户发送抢庄的用户
		serveCenter.SendTableMessage(room,message)
	case 500:
		room := common.Rooms500_.Room[tableId]
		if room.TableStatus==2 && room.Playing == 2{
			//已经有用户发送了倒计时结束请求return
			return
		}
		//设置抢庄状态为true
		Banker :=getBankerPerson(room)
		room.BankerUid = Banker["Uid"].(string)
		room.Playing = 2
		common.Rooms500_.Room[tableId] = room
		message["UserList"]=getTableUserList(common.Rooms500_.Room[tableId])
		message["Banker"] = Banker
		//向房间每位用户发送抢庄的用户
		serveCenter.SendTableMessage(room,message)
	case 1000:
		room := common.Rooms1000_.Room[tableId]
		if room.TableStatus==2 && room.Playing == 2{
			//已经有用户发送了倒计时结束请求return
			return
		}
		//设置抢庄状态为true
		Banker :=getBankerPerson(room)
		room.BankerUid = Banker["Uid"].(string)
		room.Playing = 2
		common.Rooms1000_.Room[tableId] = room
		message["UserList"]=getTableUserList(common.Rooms1000_.Room[tableId])
		message["Banker"] = Banker
		//向房间每位用户发送抢庄的用户
		serveCenter.SendTableMessage(room,message)
	}
}


//获取庄
func getBankerPerson(room common.Room)(map[string]interface{}){
	Banker:=make(map[string]interface{})
	pnum:=len(room.UserList)
	for i := 0; i < pnum-1; i++ {
		for j := 1; j < pnum-i; j++ {
			if room.UserList[j-1].RobBankerDouble >room.UserList[j].RobBankerDouble {
				tmp := room.UserList[j-1]
				room.UserList[j-1] = room.UserList[j]
				room.UserList[j] = tmp
			}
		}

	}
	maxBankerUid := room.UserList[pnum-1].Uid
	Banker["UserId"] = common.Manager.UserMapInfo[maxBankerUid].ID
	Banker["ChairId"] = room.UserList[pnum-1].ChairId
	Banker["UserName"]=common.Manager.UserMapInfo[maxBankerUid].UserName
	//json_["UserClass"]=common.Manager.UserMapInfo[uid].UserClass
	Banker["UserSrc"]=common.Manager.UserMapInfo[maxBankerUid].UserSrc
	//json_["UserSex"]=common.Manager.UserMapInfo[uid].UserSex
	Banker["UserMoney"] = common.Manager.UserMapInfo[maxBankerUid].UserActivityBalance
	Banker["Uid"] = maxBankerUid
	return Banker

}
