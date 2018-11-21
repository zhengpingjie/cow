package serveGame

import (
	"../common"
	"../serveCenter"
	"sync"
)

type BetOver struct {
	mutex sync.Mutex
}

func BetOverRest()*BetOver{
	return &BetOver{}
}

func(this *BetOver)BetOver(uid string){
	this.mutex.Lock()
	defer func() {
		this.mutex.Unlock()
		common.Wg.Done()
	}()
	tableId := common.UidJoinRoom_.UjR[uid]
	roomType := common.UidJoinRoom_.Urd[uid]
	message:=make(map[string]interface{})
	message["Action"] ="SendPokerStart"
	message["Msg"] ="发牌开始阶段"
	message["TimeOut"] =common.PokerCountDown
	switch roomType {
	case 50:
		room := common.Rooms50_.Room[tableId]
		if room.TableStatus==2 && room.Playing == 3{
			//已经有用户发送了倒计时结束请求return
			return
		}
		//游戏进入发牌阶段
		room.Playing=3
		common.Rooms50_.Room[tableId] = room
		//向房间每位用户发送抢庄的用户
		serveCenter.SendTableMessage(room,message)
	case 250:
		room := common.Rooms250_.Room[tableId]
		if room.TableStatus==2 && room.Playing == 3{
			//已经有用户发送了倒计时结束请求return
			return
		}
		//游戏进入发牌阶段
		room.Playing=3
		common.Rooms250_.Room[tableId] = room
		//向房间每位用户发送抢庄的用户
		serveCenter.SendTableMessage(room,message)
	case 500:
		room := common.Rooms500_.Room[tableId]
		if room.TableStatus==2 && room.Playing == 3{
			//已经有用户发送了倒计时结束请求return
			return
		}
		//游戏进入发牌阶段
		room.Playing=3
		common.Rooms500_.Room[tableId] = room
		//向房间每位用户发送抢庄的用户
		serveCenter.SendTableMessage(room,message)
	case 1000:
		room := common.Rooms1000_.Room[tableId]
		if room.TableStatus==2 && room.Playing == 3{
			//已经有用户发送了倒计时结束请求return
			return
		}
		//游戏进入发牌阶段
		room.Playing=3
		common.Rooms1000_.Room[tableId] = room
		//向房间每位用户发送抢庄的用户
		serveCenter.SendTableMessage(room,message)
	}

}


