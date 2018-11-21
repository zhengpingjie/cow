package common

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

//md5加密
func GetMd5(datetime string,sign string)string{
	datastr :=datetime+sign
	data:=[]byte(datastr)
	has :=md5.Sum(data)
	md5str := fmt.Sprintf("%x",has)
	return md5str
}

//验证token是否正确
func IsToken(uid string,token string)bool{
	md5str := GetMd5(Token.DateTime[uid],Token.Sign[uid])
	if token != md5str{
		//关闭连接
		Manager.Clients[uid].Close()
		//删除所有信息
		delete(Manager.Clients,uid)
		delete(Manager.UserMapInfo,uid)
		delete(Manager.RequestData,uid)
		return false
	}
	return true
}


//倒计时设置
func CountTime(num int32)int32{
	time.Sleep(1)
	if num > 0 {
		fmt.Println(num)
		CountTime(num - 1)
	}

	return num
}

//定时器
var timerPool = sync.Pool{
	New:func()interface{}{
		return time.NewTimer(time.Second)
	},
}

func SleepWithPoolTimer(d time.Duration,num int32){
	timer := timerPool.Get().(*time.Timer)
	timer.Reset(d)
	//等待触发的信号
	<- timer.C
	num--
	timerPool.Put(timer)
}

//获取用户椅子编号
func GetChairId(room Room,uid string)int32{
	var ChairId int32
	for _,roomUser:= range room.UserList{
		if uid ==roomUser.Uid{
			ChairId = roomUser.ChairId
			break
		}
	}
	return ChairId
}


//生成随机牌号 生成cnt个【start,end】不重复的随机数
func GetRandNumber(start,end,count int)[]int{

	if end<start ||(end-start) < count{
		return nil
	}

	//存储结果的slice
	nums:=make([]int,0)

	//随机数生成器
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for len(nums) < count{
		//生成随机秒数
		num := r.Intn(end-start)+start

		//查重
		exist :=false
		for _,v := range nums{
			if v == num{
				exist = true
				break
			}
		}

		if !exist{
			nums = append(nums,num)
		}
	}

	return nums
}


//获取切片中最大值
func GetSliceMaxVal(sliceArr []int)int{
	pnum:=len(sliceArr)
	for i := 0; i < pnum-1; i++ {
		for j := 1; j < pnum-i; j++ {
			if sliceArr[j-1] >sliceArr[j] {
				tmp := sliceArr[j-1]
				sliceArr[j-1] = sliceArr[j]
				sliceArr[j] = tmp
			}
		}

	}
	return sliceArr[pnum-1]

}


//检查错误
func CheckErr(err error){
	if err !=nil{
		panic(err)
	}
}

var Wg sync.WaitGroup

//移除切片中的值
func RemoveSlice(s []string,i int)[]string{
	return append(s[:i],s[i+1:]...)
}