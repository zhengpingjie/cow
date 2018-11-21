package serveGame

import (
	"../common"
	"../serveCenter"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
)

//var mutex sync.Mutex
type SelectRound struct {
	round int32
	userId int32
	mutex sync.Mutex
}

//注册选择场次
func SelectRoundRest(uid string)*SelectRound{
	data := common.Manager.RequestData[uid]
	return &SelectRound{
		round:int32(data.Data["Round"].(float64)),
		userId:int32(data.Data["UserId"].(float64)),
	}
}

//选择场次业务逻辑 返回值success 【0失败，1成功】
func (this *SelectRound)UserSelectRound(uid string){
	this.mutex.Lock()
	defer func() {
		this.mutex.Unlock()
		common.Wg.Done()
	}()
	//场次50元1元低分  250元5元低分 500元10元低分 1000元20元低分
	round:=this.round
	var UserSafeBalance float64 //用户保险柜余额
	var UserActivityBalance float64 //用户活动余额

	//判断map是否存在
	if _,ok:= common.Manager.UserMapInfo[uid];ok{
		UserSafeBalance =common.Manager.UserMapInfo[uid].UserSafeBalance
		UserActivityBalance =common.Manager.UserMapInfo[uid].UserActivityBalance
	}else{
		row:=DB.QueryRow("select UserSafeBalance,UserActivityBalance from pre_user where ID = ?",this.userId)
		row.Scan(&UserSafeBalance,&UserActivityBalance)
	}


	if round > int32(UserSafeBalance){
		common.Response.Success = "0"
		common.Response.Uid = uid
		common.Response.Action="UserPay"
		common.Response.Msg="钱包余额不足,请充值"
		common.Response.Token = common.Manager.RequestData[uid].Token
		message,_:=json.Marshal(common.Response)
		common.Manager.Clients[uid].WriteMessage(websocket.TextMessage,message)
	}else{
	//	分配桌子并创建桌子
		switch round{
		case 50:
			//如果用户更换场次 清空以前场次的数据

			if tableId,ok:=common.UidJoinRoom_.Urd[uid];ok && tableId!= 50{
				//用户为观战者
				if common.UidJoinRoom_.UserStatus[uid]==0{
					common.Rooms50_.Room[common.UidJoinRoom_.UjR[uid]]=serveCenter.ClearRoomClients(uid,common.Rooms50_.Room[common.UidJoinRoom_.UjR[uid]])
				}else if common.UidJoinRoom_.UserStatus[uid]==1 && common.Rooms50_.Room[common.UidJoinRoom_.UjR[uid]].TableStatus >= 1 {
					//用户为参战者 并且游戏已经开始 不可以随意切换房间 必须等游戏结束 注：为防止用户作弊 不可强行清除用户数据
					common.Manager.Clients[uid].WriteMessage(websocket.TextMessage,[]byte{})
					return
				}else if common.UidJoinRoom_.UserStatus[uid]==1 && common.Rooms50_.Room[common.UidJoinRoom_.UjR[uid]].TableStatus < 1{
					room:=serveCenter.ClearRoomClients(uid,common.Rooms50_.Room[common.UidJoinRoom_.UjR[uid]])
					room.UserActivity = room.UserActivity-1
					common.Rooms50_.Room[common.UidJoinRoom_.UjR[uid]] = room
				}

			}

			this.Table50(uid)
			//发送消息给房间每个用户
			tableId := common.UidJoinRoom_.UjR[uid]
			RoomInfo := selectMapMsg(common.Rooms50_.Room[tableId])
			serveCenter.SendTableMessage(common.Rooms50_.Room[tableId],RoomInfo)
			return
		case 250:
			this.Table250(uid)
			tableId := common.UidJoinRoom_.UjR[uid]
			RoomInfo := selectMapMsg(common.Rooms250_.Room[tableId])
			serveCenter.SendTableMessage(common.Rooms250_.Room[tableId],RoomInfo)
			return
		case 500:
			this.Table500(uid)
			tableId := common.UidJoinRoom_.UjR[uid]
			RoomInfo := selectMapMsg(common.Rooms500_.Room[tableId])
			serveCenter.SendTableMessage(common.Rooms500_.Room[tableId],RoomInfo)
			return
		case 1000:
			this.Table1000(uid)
			tableId := common.UidJoinRoom_.UjR[uid]
			RoomInfo := selectMapMsg(common.Rooms1000_.Room[tableId])
			serveCenter.SendTableMessage(common.Rooms1000_.Room[tableId],RoomInfo)
			return
		default:
			break

		}
	}

}


func selectMapMsg(room common.Room)map[string]interface{}{
	var RoomInfo = make(map[string]interface{})
	//获取用户列表
	RoomInfo["TableStatus"] = room.TableStatus
	if RoomInfo["TableStatus"] == 0{
		RoomInfo["Action"] = "SelectRound"
		RoomInfo["Msg"] = "继续等待其他玩家加入,等待抢庄倒计时开始"
	}
	if RoomInfo["TableStatus"] == 1{
		RoomInfo["TimeOut"] = room.RobCountDownTime
		RoomInfo["Action"] = "RoomCountDownTimeOver"
		RoomInfo["Msg"] = "游戏抢庄开始倒计时"
	}

	//RoomInfo["TableId"] = room.TableId
	//RoomInfo["Playing"] = room.Playing
	RoomInfo["UserNum"] = room.UserNum
	RoomInfo["UserActivity"] = room.UserActivity
	RoomInfo["Userlist"] = serveCenter.GetTableUserList(room)
	return  RoomInfo
}


//创建50元的桌子
func(this *SelectRound)Table50(uid string){
	//判断是否有桌子【无桌创建桌子】
	tableNum := len(common.Rooms50_.Room)
	if tableNum>0{
		var cnt = 0
		for key,val:=range common.Rooms50_.Room{
			if  val.UserNum < common.MaxNumClient {
				//用户已经存在
				fmt.Printf("%+v",val.Uids)
				if _,ok:=val.Uids[uid];ok==true {
					break
				}

				//用户加入桌子
					Room_ := common.Room{Uids:make(map[string]bool)}
					GameUser_ :=common.GameUser{}
					//人数加1
					fmt.Printf("selectRound105------%v\r\n",key)
					fmt.Printf("selectRound106------%v\r\n",val)
					val.UserNum = val.UserNum + 1
					//椅子编号为房间人数
					GameUser_.Uid = uid
					//设置默认值
					GameUser_.RobBankerDouble = 1
					GameUser_.Bet = 5
					//游戏抢庄倒计时开始
					if val.UserNum == common.MinNumClient{
						if val.TableStatus == 0{
							//进入游戏开始倒计时
							val.UserActivity =val.UserActivity+1
							val.TableStatus = 1
							GameUser_.ChairId = val.UserActivity
							//抢庄开始倒计时
							val.RobCountDownTime = common.RobCountDown
							//设置玩家为游戏参与者
							GameUser_.UserStatus = 1
							val.Playing = 1
						}
					}

					//游戏已经进入抢庄倒计时
					if val.UserNum > common.MinNumClient && val.UserNum <= common.MaxNumClient {
						//游戏未开始设置
						if  val.TableStatus == 1 {
							//设置玩家为游戏参与者
							GameUser_.UserStatus = 1
							val.UserActivity =val.UserActivity+1
							GameUser_.ChairId = val.UserActivity

						}

						//游戏已经开始
						if val.TableStatus == 2 {
							//设置玩家为游戏参与者
							GameUser_.UserStatus = 0
						}

					}


					//用户已经加入

					Room_.UserNum=val.UserNum
					Room_.UserActivity = val.UserActivity
					Room_.Uids[uid] = true
					Room_.TableStatus = val.TableStatus
					Room_.Playing = val.Playing
					Room_.RobCountDownTime = val.RobCountDownTime
					Room_.TableId = key
					//添加用户到房间
					if len(val.UserList)>0{
						for _,user:=range val.UserList{
							Room_.UserList = append(Room_.UserList, user)
							Room_.Uids[user.Uid] = true
						}

					}

					Room_.UserList = append(Room_.UserList, GameUser_)
					common.Rooms50_.Room[key] = Room_
					common.UidJoinRoom_.UjR[uid] = val.TableId
					common.UidJoinRoom_.Urd[uid] = 50
					common.UidJoinRoom_.UserStatus[uid] = GameUser_.UserStatus

					fmt.Printf("selectRound147-----%+v\r\n",common.Rooms50_.Room[key])
					break


			}

			cnt++

		}

		//	桌子已经满了 创建新桌子
		if cnt >= tableNum {
			//无桌子创建桌子
			tableId := serveCenter.Createtable50(tableNum+1)
			common.UidJoinRoom_.UjR[uid] = tableId
			common.UidJoinRoom_.Urd[uid] = 50
			common.UidJoinRoom_.UserStatus[uid] = 1

			//添加用户到桌子
			serveCenter.AddUserToTable50(uid)
		}

	}else{
		//无桌子创建桌子
		tableId := serveCenter.Createtable50(0)
		fmt.Printf("selectRound170-----------%+v\r\n",common.Rooms50_.Room)
		common.UidJoinRoom_.UjR[uid] = tableId
		common.UidJoinRoom_.Urd[uid] = 50
		common.UidJoinRoom_.UserStatus[uid] = 1
		//添加用户到桌子
		serveCenter.AddUserToTable50(uid)
		fmt.Printf("selectRound174-----------%+v\r\n",common.Rooms50_.Room)
	}

}


//创建250的桌子
func(this *SelectRound)Table250(uid string){
	//判断是否有桌子【无桌创建桌子】
	tableNum := len(common.Rooms250_.Room)
	if tableNum>0{
		var cnt = 0
		for key,val:=range common.Rooms250_.Room{
			if _,ok:=val.Uids[uid];!ok && val.UserNum < common.MaxNumClient {
				Room_ := common.Room{Uids:make(map[string]bool)}
				//用户加入桌子
				//人数加1
				GameUser_ := common.GameUser{}
				fmt.Printf("%v",key)
				fmt.Printf("%v",val)
				val.UserNum = val.UserNum + 1
				//椅子编号为房间人数
				GameUser_.ChairId = val.UserNum
				GameUser_.Uid = uid

				if val.UserNum >= common.MinNumClient && val.UserNum <= 5 {
					//游戏未开始设置
					if val.TableStatus == 0 || val.TableStatus == 1 {
						//进入游戏开始倒计时
						val.TableStatus = 1
						//设置玩家为游戏参与者
						GameUser_.UserStatus = 1
					}

					//游戏未开始
					if val.TableStatus == 2 {
						//设置玩家为游戏参与者
						GameUser_.UserStatus = 0
					}

				}


				//用户已经加入

				Room_.UserNum=val.UserNum
				Room_.Uids[uid] = true
				Room_.TableStatus = val.TableStatus

				//添加用户到房间
				Room_.UserList = append(Room_.UserList, GameUser_)
				common.Rooms250_.Room[key] = Room_
				common.UidJoinRoom_.UjR[uid] = val.TableId

				fmt.Printf("%v",common.Rooms250_.Room[key])
				break


			}


			cnt++

		}

		//	桌子已经满了 创建新桌子
		if cnt >= tableNum {
			//无桌子创建桌子
			common.UidJoinRoom_.UjR[uid] = serveCenter.Createtable250(tableNum+1)

			common.Rooms250_.Room[common.UidJoinRoom_.UjR[uid]] = common.Rooms250_.Room[common.UidJoinRoom_.UjR[uid]]

			//添加用户到桌子
			serveCenter.AddUserToTable250(uid)
		}

	}else{
		//无桌子创建桌子
		common.UidJoinRoom_.UjR[uid] = serveCenter.Createtable250(0)
		common.Rooms250_.Room[common.UidJoinRoom_.UjR[uid]] = common.Rooms250_.Room[common.UidJoinRoom_.UjR[uid]]

		//添加用户到桌子
		serveCenter.AddUserToTable250(uid)
	}

}

//创建500元的桌子
func(this *SelectRound)Table500(uid string){
	//判断是否有桌子【无桌创建桌子】
	tableNum := len(common.Rooms500_.Room)
	if tableNum>0{
		var cnt = 0
		for key,val:=range common.Rooms500_.Room{
			if _,ok:=val.Uids[uid];!ok && val.UserNum < common.MaxNumClient {
				Room_ := common.Room{Uids:make(map[string]bool)}
				//用户加入桌子
				GameUser_ := common.GameUser{}
				//人数加1
				fmt.Printf("%v",key)
				fmt.Printf("%v",val)
				val.UserNum = val.UserNum + 1
				//椅子编号为房间人数
				GameUser_.ChairId = val.UserNum
				GameUser_.Uid = uid

				if val.UserNum >= common.MinNumClient && val.UserNum <= 5 {
					//游戏未开始设置
					if val.TableStatus == 0 || val.TableStatus == 1 {
						//进入游戏开始倒计时
						//设置玩家为游戏参与者
						GameUser_.UserStatus = 1
					}

					//游戏未开始
					if val.TableStatus == 2 {
						//设置玩家为游戏参与者
						GameUser_.UserStatus = 0
					}

				}


				//用户已经加入
				Room_.UserNum=val.UserNum
				Room_.Uids[uid] = true
				Room_.TableStatus = val.TableStatus

				//添加用户到房间
				Room_.UserList = append(Room_.UserList, GameUser_)
				common.Rooms500_.Room[key] = Room_
				common.UidJoinRoom_.UjR[uid] = val.TableId

				fmt.Printf("%v",common.Rooms500_.Room[key])
				break


			}


			cnt++

		}

		//	桌子已经满了 创建新桌子
		if cnt >= tableNum {
			//无桌子创建桌子
			common.UidJoinRoom_.UjR[uid] = serveCenter.Createtable500(tableNum+1)

			common.Rooms500_.Room[common.UidJoinRoom_.UjR[uid]] = common.Rooms500_.Room[common.UidJoinRoom_.UjR[uid]]

			//添加用户到桌子
			serveCenter.AddUserToTable500(uid)
		}

	}else{
		//无桌子创建桌子
		common.UidJoinRoom_.UjR[uid] = serveCenter.Createtable500(0)
		common.Rooms500_.Room[common.UidJoinRoom_.UjR[uid]] = common.Rooms500_.Room[common.UidJoinRoom_.UjR[uid]]

		//添加用户到桌子
		serveCenter.AddUserToTable500(uid)
	}


}
//创建1000元的桌子
func(this *SelectRound)Table1000(uid string){
	//判断是否有桌子【无桌创建桌子】
	tableNum := len(common.Rooms1000_.Room)
	if tableNum>0{
		var cnt = 0
		for key,val:=range common.Rooms1000_.Room{
			if _,ok:=val.Uids[uid];!ok && val.UserNum < common.MaxNumClient {
				Room_ := common.Room{Uids:make(map[string]bool)}
				GameUser_ :=common.GameUser{}
				//用户加入桌子
				//人数加1
				fmt.Printf("%v",key)
				fmt.Printf("%v",val)
				val.UserNum = val.UserNum + 1
				//椅子编号为房间人数
				GameUser_.ChairId = val.UserNum
				GameUser_.Uid = uid

				if val.UserNum >= common.MinNumClient && val.UserNum <= 5 {
					//游戏未开始设置
					if val.TableStatus == 0 || val.TableStatus == 1 {
						//进入游戏开始倒计时
						val.TableStatus = 1
						//设置玩家为游戏参与者
						GameUser_.UserStatus = 1
					}

					//游戏未开始
					if val.TableStatus == 2 {
						//设置玩家为游戏参与者
						GameUser_.UserStatus = 0
					}

				}


				//用户已经加入

				Room_.UserNum=val.UserNum
				Room_.Uids[uid] = true
				Room_.TableStatus = val.TableStatus

				//添加用户到房间
				Room_.UserList = append(Room_.UserList, GameUser_)
				common.Rooms1000_.Room[key] = Room_
				common.UidJoinRoom_.UjR[uid] = val.TableId

				fmt.Printf("%v",common.Rooms1000_.Room[key])
				break


			}


			cnt++

		}

		//	桌子已经满了 创建新桌子
		if cnt >= tableNum {
			//无桌子创建桌子
			common.UidJoinRoom_.UjR[uid] = serveCenter.Createtable1000(tableNum+1)

			common.Rooms1000_.Room[common.UidJoinRoom_.UjR[uid]] = common.Rooms1000_.Room[common.UidJoinRoom_.UjR[uid]]

			//添加用户到桌子
			serveCenter.AddUserToTable1000(uid)
		}

	}else{
		//无桌子创建桌子
		common.UidJoinRoom_.UjR[uid] = serveCenter.Createtable1000(0)
		common.Rooms1000_.Room[common.UidJoinRoom_.UjR[uid]] = common.Rooms1000_.Room[common.UidJoinRoom_.UjR[uid]]

		//添加用户到桌子
		serveCenter.AddUserToTable1000(uid)
	}

}