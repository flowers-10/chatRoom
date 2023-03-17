package process2

import (
	"encoding/json"
	"fmt"
	"gocode/chatroom/common/message"
	"gocode/chatroom/server/model"
	"gocode/chatroom/server/utils"
	"net"
)

type UserProcess struct {
	//字段?
	Conn net.Conn
	//该conn是哪个用户
	UserId int
}

// 通知所有在线用户的方法
// 要通知其他用户上线
func (u *UserProcess) NotifyOthersOnlineUser(userId int) {
	// 遍历 onLineUsers，然后一个一个的发送 NotifyUserStatusMes
	for id, up := range userMgr.onlineUsers {
		if id == userId {
			continue
		}
		// 开始通知
		up.NotifyMeOnline(userId)
	}

}
func (u *UserProcess) NotifyMeOnline(userId int) {
	// 组装NotifyUserStatusMes
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var _notifyUserStatusMes message.NotifyUserStatusMes
	_notifyUserStatusMes.UserId = userId
	_notifyUserStatusMes.Status = message.UserOnline

	// 将_notifyUserStatusMes序列化
	data, err := json.Marshal(_notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	// 序列化的mes赋值给mes.Data
	mes.Data = string(data)
	// 序列化，发送
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	// 发送,创建transfer实例
	tf := utils.Transfer{
		Conn: u.Conn,
	}

	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyMeOnline err = ", err)
		return
	}
}

// serverProcessLogin,处理登录请求
func (u *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
	// 核心代码
	// 1先从mes 中取出 mes.Data,并反序列化成LoginMes
	var _LoginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &_LoginMes)
	if err != nil {
		fmt.Println("json.Unmarshal err", err)
		return
	}
	// 2声明resMes和LoginResMes赋值
	var resMes message.Message
	resMes.Type = message.LoginResMesType
	var _LoginResMes message.LoginResMes
	// 我们需要到redis数据库去完成登录验证
	//1.使用model.MyUserDao 到redis去验证
	user, err := model.MyUserDao.Login(_LoginMes.UserId, _LoginMes.UserPwd)
	if err != nil {

		if err == model.ERROR_USER_NOTEXISTS {
			_LoginResMes.Code = 500 //不合法
			_LoginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			_LoginResMes.Code = 403 //不合法
			_LoginResMes.Error = err.Error()
		} else {
			_LoginResMes.Code = 505 //不合法
			_LoginResMes.Error = "服务器内部错误.."
		}

	} else {
		// 合法
		_LoginResMes.Code = 200
		// 因为用户登录成功，把该登录成功的用户放到userMgr中
		u.UserId = _LoginMes.UserId
		// 把登录成功的用户收集到onlineUsers
		userMgr.AddOnlineUsers(u)
		u.NotifyOthersOnlineUser(_LoginMes.UserId)
		// 便利 userMgr.onlineUsers获取所有在线用户
		for id, _ := range userMgr.onlineUsers {
			_LoginResMes.UsersId = append(_LoginResMes.UsersId, id)
		}
		fmt.Println(user, "登录成功")
	}

	// // 如果用户id =100，密码123456认为合法否则不合法
	// if _LoginMes.UserId == 100 && _LoginMes.UserPwd == "123456" {
	// 	// 合法
	// 	_LoginResMes.Code = 200
	// } else {
	// 	_LoginResMes.Code = 500 //不合法
	// 	_LoginResMes.Error = "该用户不存在，请注册使用"
	// }
	// 3将 loginResMes序列化
	data, err := json.Marshal(_LoginResMes)
	if err != nil {
		fmt.Println("json.Marshal err", err)
		return
	}
	// 4.将data 赋值给 resMes
	resMes.Data = string(data)
	// 5.对resMes序列化准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("jon.Marshal fail", err)
		return
	}
	// 6.发送data
	// 因为使用了分层模式（mvc），我们先创建一个Transfer实例，然后读取
	tf := &utils.Transfer{
		Conn: u.Conn,
	}
	tf.WritePkg(data)
	return
}

// 处理注册
func (u *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	// 核心代码

	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal err", err)
		return
	}

	var resMes message.Message
	resMes.Type = message.RegisterResMesType
	var registerResMes message.RegisterResMes
	// 我们需要到redis数据库去完成登录验证
	//1.使用model.MyUserDao 到redis去验证
	err = model.MyUserDao.Register(&registerMes.User)
	if err != nil {

		if err == model.ERROR_USER_EXISTS {
			registerResMes.Code = 505 //不合法
			registerResMes.Error = model.ERROR_USER_EXISTS.Error()
		} else {
			registerResMes.Code = 505 //不合法
			registerResMes.Error = "注册发生未知错误.."
		}

	} else {
		// 合法
		registerResMes.Code = 200
		fmt.Println("注册成功")
	}

	// 3序列化
	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json.Marshal err", err)
		return
	}
	// 4.将data 赋值给 resMes
	resMes.Data = string(data)
	// 5.对resMes序列化准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("jon.Marshal fail", err)
		return
	}
	// 6.发送data
	// 因为使用了分层模式（mvc），我们先创建一个Transfer实例，然后读取
	tf := &utils.Transfer{
		Conn: u.Conn,
	}
	tf.WritePkg(data)
	return
}
