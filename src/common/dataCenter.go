package common

import "github.com/gorilla/websocket"

type ClientManager struct {
	Clients    map[string]*websocket.Conn
	UserMapInfo map[string]UserLoginInfo
	RequestData map[string]RequestParam
}


var Manager = ClientManager{
	Clients:  make(map[string]*websocket.Conn),
	UserMapInfo:make(map[string]UserLoginInfo),
	RequestData:make(map[string]RequestParam),
}


//用户登陆信息
type UserLoginInfo struct {
	ID int
	UserName string
	UserPhone string
	UserPassword string
	UserBirthday string
	UserClass int               //用户分组/用户等级
	UserSrc  string             //用户头像
	UserSex int			        //用户性别
	UserSafeBalance  float64    //用户保险柜余额
	UserActivityBalance float64 //用户活动余额
}

var UserInfo = UserLoginInfo{}



//请求数据
type RequestParam struct {
	Uid string `json:",omitempty"`
	Token string `json:",omitempty"` //登陆成功返回的token
	Action string `json:",omitempty"` //触发的动作
	Data map[string]interface{} `json:",omitempty"`
}

//返回数据
type ResponseParam struct {
	Success string  //0=>失败 1=>成功
	Msg  string //消息提示
	OnlineNum int   //在线人事
	Uid string `json:",omitempty"`
	Token string `json:",omitempty"` //登陆成功返回的token
	Action string `json:",omitempty"` //触发的动作
	Data interface{} `json:",omitempty"`
}

var Response = ResponseParam{}

//用户签名和时间戳
type UserToken struct {
	DateTime  map[string]string  //时间戳
	Sign     map[string]string  //用户签名
}

var Token = UserToken{DateTime:make(map[string]string), Sign : make(map[string]string)}





//入场限制50元 1元的低分
type Rooms50 struct {
	Room map[string]Room
}

//入场限制250元 5元的低分
type Rooms250 struct {
	Room map[string]Room
}

//入场限制500元 10元的低分

type Rooms500 struct {
	Room map[string]Room
}



//入场限制1000元 20元的低分
type Rooms1000 struct {
	Room map[string]Room
}


type Room struct {
	TableId string
	Uids    map[string]bool
	TableStatus int  //[0=>等待，1=>开始倒计时 2=>开始]房间状态
	UserList []GameUser
	UserNum int32  //房间总人数
	UserActivity int32 //房间活动人数
	Playing int32  //游戏中状态【1=>抢庄中 2=>押注中 3=>发牌 4=>摊牌 5=>结算中】
	RobCountDownTime int32     //倒计时时间
	CartList []int  //所有牌string
	BankerUid  string//房间庄uid
}

//游戏中的用户
type GameUser struct{
	Uid      string    // 【编号】
	RobBankerDouble int     //【0=>不抢，1=>1倍,2=>2倍 ,4=>4倍】
	Bet int  //【5=>5倍,10=>10倍,15=>15倍,20=>20倍】
	IsBet bool //是否下注
	Isbanker bool //是否抢庄
	IsbShowPoker bool // 是否摊牌
	ChairId int32 //【座位号】
	UserStatus int32 //【0=>观察者 1=>玩家】
	UserCartList []int //玩家牌
}

var Rooms50_ = Rooms50{make(map[string]Room)}
var Rooms250_ = Rooms250{make(map[string]Room)}
var Rooms500_ = Rooms500{make(map[string]Room)}
var Rooms1000_ = Rooms1000{make(map[string]Room)}
//var Room_ = Room{Uids:make(map[string]bool)}

//用户关联的房间和场次
type UidJoinRoom struct {
	UjR map[string]string
	Urd map[string]int32
	UserStatus map[string] int32
}
var  UidJoinRoom_ = UidJoinRoom{make(map[string]string),make(map[string]int32),make(map[string] int32)}


//记录房间用户比牌记录tableid uid

type RoomThenRecord struct {
	RoomRecord map[string]map[string]UserRecord
}
var RoomThenRecord_ = RoomThenRecord{RoomRecord:make(map[string]map[string]UserRecord)}

//用户牌记录
type UserRecord struct {
	Uid string
	CowNum int //[0=>无牛，1=>牛1.....9=>牛9,10=>牛牛]
	MaxPoker int //最大牌
	MaxPokerColor int //【4321 分别代表黑桃，红桃，梅花，方块】
	PokersCount  int //所有牌的数字总和
	Poker1Num int//牌1牌号
	Poker1Color int //【4321 分别代表黑桃，红桃，梅花，方块】
	Poker2Num int //牌2牌号
	Poker2Color int //【4321 分别代表黑桃，红桃，梅花，方块】
	Poker3Num int //牌3牌号
	Poker3Color int //【4321 分别代表黑桃，红桃，梅花，方块】
	Poker4Num int //牌4牌号
	Poker4Color int //【4321 分别代表黑桃，红桃，梅花，方块】
	Poker5Num int //牌5牌号
	Poker5Color int //【4321 分别代表黑桃，红桃，梅花，方块】
}