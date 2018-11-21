package serveGame

import (
	"../common"
	"../serveCenter"
)

type RobBankerReady struct {
	chairId int32
}

func RobBankerReadyRest(uid string)*RobBankerReady{
	data := common.Manager.RequestData[uid]
	return &RobBankerReady{
		chairId:data.Data["chairId"].(int32),
	}
}


func (this *RobBankerReady)RobBankerReady(uid string){
	message:=make(map[string]interface{})
	message["Content"]="抢庄开始"
	message["Action"]="RobStart"
	serveCenter.SendClientMessage(uid,message)
}

