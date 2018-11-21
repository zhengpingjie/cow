package serveGame

import (
	"../common"
	"sync"
)

type CountDownTime struct {
	timeout int32
	mutex sync.Mutex
}

//注册struct
func CountDownTimeRest(uid string)*CountDownTime{
	data := common.Manager.RequestData[uid]
	return &CountDownTime{
		timeout:int32(data.Data["TimeOut"].(float64)),
	}
}


func(this *CountDownTime)UpdateCountDownTime(uid string){
	this.mutex.Lock()
	defer func() {
		this.mutex.Unlock()
		common.Wg.Done()
	}()
	tableId:=common.UidJoinRoom_.UjR[uid]
	roomType := common.UidJoinRoom_.Urd[uid]
	switch roomType {
	case 50:
		room := common.Rooms50_.Room[tableId]
		room.RobCountDownTime= this.timeout
		common.Rooms50_.Room[tableId] = room
	case 250:
		room := common.Rooms50_.Room[tableId]
		room.RobCountDownTime= this.timeout
		common.Rooms50_.Room[tableId] = room
	case 500:
		room := common.Rooms50_.Room[tableId]
		room.RobCountDownTime= this.timeout
		common.Rooms50_.Room[tableId] = room
	case 1000:
		room := common.Rooms50_.Room[tableId]
		room.RobCountDownTime= this.timeout
		common.Rooms50_.Room[tableId] = room

	}
}
