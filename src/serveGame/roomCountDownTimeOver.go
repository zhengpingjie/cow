package serveGame

import (
	"../common"
	"../serveCenter"
	"sync"
)

type CountDownTimeOver struct {
	mutex sync.Mutex
}

func TimeOverRest()*CountDownTimeOver{
	return &CountDownTimeOver{
		}
}


func (this *CountDownTimeOver)TimeOver(uid string){
	this.mutex.Lock()
	defer func() {
		this.mutex.Unlock()
		common.Wg.Done()
	}()

	tableId := common.UidJoinRoom_.UjR[uid]
	if common.Rooms50_.Room[tableId].TableStatus==2 && common.Rooms50_.Room[tableId].Playing == 1{
		//已经有用户发送了倒计时结束请求return
		return
	}
	roomType := common.UidJoinRoom_.Urd[uid]
	//获取用户id以及椅子编号
	message:=make(map[string]interface{})
	var room common.Room
	switch roomType {
	case 50:
		room := common.Rooms50_.Room[tableId]
		if room.TableStatus==2 && room.Playing == 1{
			//已经有用户发送了倒计时结束请求return
			return
		}
		//修改桌面状态
		room.TableStatus = 2
		room.Playing = 1
		common.Rooms50_.Room[tableId] = room

	case 250:
		room := common.Rooms50_.Room[tableId]
		if room.TableStatus==2 && room.Playing == 1{
			//已经有用户发送了倒计时结束请求return
			return
		}
		//修改桌面状态
		room.TableStatus = 2
		room.Playing = 1
		common.Rooms250_.Room[tableId] = room

	case 500:
		room := common.Rooms50_.Room[tableId]
		if room.TableStatus==2 && room.Playing == 1{
			//已经有用户发送了倒计时结束请求return
			return
		}
		//修改桌面状态
		room.TableStatus = 2
		room.Playing = 1
		common.Rooms500_.Room[tableId] = room
	case 1000:
		room := common.Rooms50_.Room[tableId]
		if room.TableStatus==2 && room.Playing == 1{
			//已经有用户发送了倒计时结束请求return
			return
		}
		//修改桌面状态
		room.TableStatus = 2
		room.Playing = 1
		common.Rooms1000_.Room[tableId] = room


	}

	message["Action"]="RobBankerStart"
	message["Msg"]="抢庄开始"
	message["TimeOut"] = common.RobBankerValidTime
	//向房间每个用户发送抢庄开始消息
	serveCenter.SendTableMessage(room,message)
}