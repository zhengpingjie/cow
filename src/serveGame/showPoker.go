package serveGame

import (
	"../common"
)

type ShowPoker struct {
	chairId int32
}

func ShowPokerRest(uid string)*ShowPoker{
	data := common.Manager.RequestData[uid]
	return &ShowPoker{
		chairId:int32(data.Data["ChairId"].(float64)),
	}
}

func (this *ShowPoker)ShowPoker(uid string){
	tableId := common.UidJoinRoom_.UjR[uid]
	roomType := common.UidJoinRoom_.Urd[uid]
	switch roomType {
	case 50:
		room := common.Rooms50_.Room[tableId]
		//整理用户牌数
		setUserPoker(room)
		//摊牌结算
		setAccount(room)


	}
}

//计算是否有牛
func setUserPoker(room common.Room){
 //整理牌号对应花牌
 setPoker(room)
 //设置最大牌数
setMaxPoker(room.TableId)
 //计算牛
setCowNum(room.TableId)


}

//摊牌结算
func setAccount(room common.Room){
	var bankerPoker common.UserRecord

	thenRecord := common.RoomThenRecord_.RoomRecord[room.TableId]
	//获取庄家牌
	for key,pokerinfo:= range thenRecord{
		if key==room.BankerUid{
			bankerPoker=pokerinfo
			delete(thenRecord,key)
		}
	}

	//与庄家比牌结算
	for _,user := range thenRecord{
		winuid,times:=thenPoker(bankerPoker,user)
		if winuid == ""{
			continue
		}
		loseuid:=user.Uid
		account(winuid,loseuid,times,room)
	}

}

//结算
func account(winuid,loseuid string,times int,room common.Room){
	accountMap:=make(map[string]map[string]interface{})
	for _,userInfo:=range room.UserList{
		if userInfo.Uid == winuid{
			income :=userInfo.Bet*userInfo.RobBankerDouble*times
			accountMap[winuid]["type"] = "add"
			accountMap[winuid]["income"] = income
			accountMap[winuid]["userid"] = common.Manager.UserMapInfo[winuid].ID

			accountMap[loseuid]["type"] = "reduce"
			accountMap[loseuid]["pay"] = income
			accountMap[loseuid]["userid"] = common.Manager.UserMapInfo[loseuid].ID

			tx,err :=DB.Begin()
			common.CheckErr(err)
			stmt1,err11:=tx.Prepare("Update pre_user set UserActivityBalance = UserActivityBalance-? where ID=?")
			common.CheckErr(err11)
			res,err12 :=stmt1.Exec(income,accountMap[loseuid]["userid"])
			common.CheckErr(err12)
			ids1,err13:=res.LastInsertId()
			common.CheckErr(err13)

			stmt2,err21:=tx.Prepare("Update pre_user set UserActivityBalance = UserActivityBalance+? where ID=?")
			common.CheckErr(err21)
			res2,err22 :=stmt2.Exec(income,accountMap[winuid]["userid"])
			common.CheckErr(err22)
			ids2,err23:=res2.LastInsertId()
			common.CheckErr(err23)
			if ids1>0 && ids2>0{
				tx.Commit()
			}else{
				tx.Rollback()
			}
			break
		}
		//数据库操作


	}




}

//比牌
func thenPoker(bankerPoker,userPoker common.UserRecord)(string,int){
	var winuid string
	var times int

	if bankerPoker == userPoker{
		winuid = ""
		return winuid,times
	}

	if bankerPoker.CowNum>0 || userPoker.CowNum>0{
		if bankerPoker.CowNum == 10 || userPoker.CowNum == 10{
			times = 3
		}else{
			times = 2
		}
	}

	//有牛
	if bankerPoker.CowNum >0{

		//庄家有牛 玩家无牛
		if userPoker.CowNum <= 0{
			return bankerPoker.Uid,times
		}

		//庄家有牛 玩家也有牛
		if bankerPoker.CowNum > userPoker.CowNum{
			winuid = bankerPoker.Uid
		}else if bankerPoker.CowNum == userPoker.CowNum{

			if bankerPoker.PokersCount >userPoker.PokersCount{
				//比牌的点数
				winuid = bankerPoker.Uid
				return winuid,times
			}
			if bankerPoker.MaxPoker >userPoker.MaxPoker{
				//比最大牌
				winuid = bankerPoker.Uid
				return winuid,times
			}
			if bankerPoker.MaxPokerColor > userPoker.MaxPokerColor{
				//比最大牌花色
				winuid = bankerPoker.Uid
				return winuid,times
			}

			winuid = userPoker.Uid
			return winuid,times

		}else{
			if userPoker.CowNum==10{
				times = 3
			}else{
				times = 2
			}
			winuid = userPoker.Uid
		}

	}

	//无牛
	if bankerPoker.CowNum <= 0{
		//庄家无牛 玩家有牛
		if userPoker.CowNum>0{
			winuid = userPoker.Uid
			return winuid,times
		}

		times = 1
		if bankerPoker.PokersCount >userPoker.PokersCount{
			//比牌的点数
			winuid = bankerPoker.Uid
			return winuid,times
		}
		if bankerPoker.MaxPoker >userPoker.MaxPoker{
			//比最大牌
			winuid = bankerPoker.Uid
			return winuid,times
		}
		if bankerPoker.MaxPokerColor > userPoker.MaxPokerColor{
			//比最大牌花色
			winuid = bankerPoker.Uid
			return winuid,times
		}

		winuid = userPoker.Uid

	}

	return winuid,times

}

//整理牌号对应花牌
func setPoker(room common.Room){
	UserRecord := common.UserRecord{}
	userList := room.UserList
	for key,user := range userList{

		for _,poker := range user.UserCartList{
			UserRecord.Uid = user.Uid
			if key == 0{
				PokerNum,PokerColor:=setPokerNum(poker)
				UserRecord.Poker1Num = PokerNum
				UserRecord.Poker1Color = PokerColor
				UserRecord.PokersCount = UserRecord.PokersCount+PokerNum
			}else if key == 1{
				PokerNum,PokerColor:=setPokerNum(poker)
				UserRecord.Poker2Num = PokerNum
				UserRecord.Poker2Color = PokerColor
				UserRecord.PokersCount = UserRecord.PokersCount+PokerNum
			}else if key == 2{
				PokerNum,PokerColor:=setPokerNum(poker)
				UserRecord.Poker3Num = PokerNum
				UserRecord.Poker3Color = PokerColor
				UserRecord.PokersCount = UserRecord.PokersCount+PokerNum
			}else if key == 3{
				PokerNum,PokerColor:=setPokerNum(poker)
				UserRecord.Poker4Num = PokerNum
				UserRecord.Poker4Color = PokerColor
				UserRecord.PokersCount = UserRecord.PokersCount+PokerNum
			}else if key == 4{
				PokerNum,PokerColor:=setPokerNum(poker)
				UserRecord.Poker5Num = PokerNum
				UserRecord.Poker5Color = PokerColor
				UserRecord.PokersCount = UserRecord.PokersCount+PokerNum
			}
		}

		common.RoomThenRecord_.RoomRecord[room.TableId][user.Uid] = UserRecord
	}


}


//设置最大牌数
func setMaxPoker(tableid string){
	userinfo:=common.RoomThenRecord_.RoomRecord[tableid]
	userPokerSlice:=[]int{}
	userPokersColor:=[]int{}
	for _,UserRecord:=range userinfo{
		userPokerSlice = append(userPokerSlice,UserRecord.Poker1Num,UserRecord.Poker2Num,UserRecord.Poker3Num,UserRecord.Poker4Num,UserRecord.Poker5Num)
		userPokersColor = append(userPokersColor,UserRecord.Poker1Color,UserRecord.Poker2Color,UserRecord.Poker3Color,UserRecord.Poker4Color,UserRecord.Poker5Color)
		MaxPoker,MaxPokerColor:=getMaxNum(userPokerSlice,userPokersColor)
		UserRecord:=common.RoomThenRecord_.RoomRecord[tableid][UserRecord.Uid]
		UserRecord.MaxPoker = MaxPoker
		UserRecord.MaxPokerColor = MaxPokerColor
		common.RoomThenRecord_.RoomRecord[tableid][UserRecord.Uid] = UserRecord
	}
}

//获取最大值
func getMaxNum(userPokers,userPokersColor []int)(int,int){
	for i := 0; i < 5; i++ {
		for j := 1; j < 5-i; j++ {
			if userPokers[j-1] >userPokers[j] {
				Poker :=userPokers[j-1]
				userPokers[j-1] = userPokers[j]
				userPokers[j] = Poker

				color := userPokersColor[j-1]
				userPokersColor[j-1] = userPokersColor[j]
				userPokersColor[j] = color
			}
		}

	}
	return userPokers[4],userPokersColor[4]

}


//设置牛
func setCowNum(tableid string){
	userinfo := common.RoomThenRecord_.RoomRecord[tableid]
	for _,UserRecord:=range userinfo{
		CowNum:=getCowNum(UserRecord)
		UserRecord.CowNum = CowNum
		common.RoomThenRecord_.RoomRecord[tableid][UserRecord.Uid] = UserRecord
	}

}

//计算牛【5中随机抽取三位为10的倍数】总共有10种组合
func getCowNum(UserRecord common.UserRecord)int{
	var cowArr []int
	var cowNum int
	if UserRecord.Poker1Num >10{
		UserRecord.Poker1Num = 10
	}

	if UserRecord.Poker2Num >10{
		UserRecord.Poker2Num =10
	}

	if UserRecord.Poker3Num >10{
		UserRecord.Poker3Num = 10
	}

	if UserRecord.Poker4Num >10{
		UserRecord.Poker4Num = 10
	}

	if UserRecord.Poker5Num >10{
		UserRecord.Poker5Num = 10
	}
	one := UserRecord.Poker2Num
	two := UserRecord.Poker2Num
	three := UserRecord.Poker3Num
	four := UserRecord.Poker4Num
	five := UserRecord.Poker5Num
	for i:=0;i<10;i++{
		//123组合
		if (one+two+three)%10 == 0{
			if (four+five)>10{
				cowArr = append(cowArr,four+five -10)
			}else{
				cowArr = append(cowArr,four+five)
			}

		}
		//124组合
		if(one+two+four)%10 ==0{
			if(three+five)>10{
				cowArr = append(cowArr,three+five -10)
			}else{
				cowArr = append(cowArr,three+five)
			}
		}
		//125组合
		if(one+two+five)%10 == 0{
			if(three+five)>10{
				cowArr = append(cowArr, three+five -10)
			}else{
				cowArr = append(cowArr, three+five)
			}
		}
		//134组合
		if(one+three+four)%10 == 0{
			if(two+five)>10{
				cowArr = append(cowArr, two+five -10)
			}else{
				cowArr = append(cowArr, two+five)
			}

		}
		//135组合
		if(one+three+five)%10 == 0{
			if(two+four)>10{
				cowArr = append(cowArr, two+four -10)
			}else{
				cowArr = append(cowArr, two+four)
			}
		}
		//145组合
		if(one+four+five)%10 == 0{
			if(two+three)>10{
				cowArr = append(cowArr, two+three -10)
			}else{
				cowArr = append(cowArr, two+three)
			}
		}

		//234组合
		if(two+three+four)%10 == 0{
			if(one+five)>10{
				cowArr = append(cowArr, one+five -10)
			}else{
				cowArr = append(cowArr, one+five)
			}
		}
		//235组合
		if(two+three+five)%10 == 0{
			if(one+four)>10{
				cowArr = append(cowArr, one+four -10)
			}else{
				cowArr = append(cowArr, one+four)
			}
		}
		//245组合
		if(two+four+five)%10 == 0{
			if(one+four)>10{
				cowArr = append(cowArr, one+three -10)
			}else{
				cowArr = append(cowArr, one+three)
			}
		}
		//345组合
		if(three+four+five)%10 == 0{
			if(one+two)>10{
				cowArr = append(cowArr, one+two -10)
			}else{
				cowArr = append(cowArr, one+two)
			}
		}



	}
	lenth := len(cowArr)
	if lenth >0{
		if lenth==1{
			cowNum = cowArr[0]
		}else{
			cowNum = common.GetSliceMaxVal(cowArr)
		}
	}

	return cowNum
}


func setPokerNum(poker int)(int,int){
	var PokerNum,PokerColor int
	if poker <=13 {
		PokerNum = poker
		PokerColor = common.PokerBlack //黑桃
	}else if poker > 13 &&  poker<=26{
		PokerNum = poker-13
		PokerColor = common.PokerRed //红桃
	}else if poker>26 && poker<=39{
		PokerNum = poker-26
		PokerColor = common.PokerFlower //梅花
	}else if poker>39 && poker<=52{
		PokerNum = poker-39
		PokerColor = common.PokerSquare //方块
	}
	return PokerNum,PokerColor
}