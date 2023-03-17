package message

const (
	LoginMesType            = "LoginMes"
	LoginResMesType         = "LoginResMes"
	RegisterMesType         = "RegisterMes"
	RegisterResMesType      = "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	SmsMesType              = "SmsMes"
)

// 定义用户状态常量
const (
	UserOnline = iota
	UserOffline
	UserBusyStatus
)

type Message struct {
	Type string `json:"type"` //消息类型
	Data string `json:"data"` //消息类型
}

// 定义两个消息。。后面需要再添加

type LoginMes struct {
	UserId   int    `json:"userId"`   //用户序号
	UserPwd  string `json:"userPwd"`  //用户密码
	UserName string `json:"UserName"` //用户名
}

type LoginResMes struct {
	Code    int    `json:"code"` //返回状态码，500表示该用户未注册，200登录成功
	UsersId []int  //在线用户id的切片
	Error   string `json:"Error"` //返回错误信息
}

type RegisterMes struct {
	User User `json:"user"` //User结构体
}

type RegisterResMes struct {
	Code  int    `json:"code"`  //返回状态码，400表示该用户已经占用，200注册成功
	Error string `json:"Error"` //返回错误信息
}

// 为了配合服务器端推送用户状态变化的消息
type NotifyUserStatusMes struct {
	UserId int `json:"userId"` //用户id
	Status int `json:"status"` //用户状态
}

// SmsMes发送信息
type SmsMes struct {
	Content string `json:"content"` //内容
	User           //匿名结构体，继承
}
