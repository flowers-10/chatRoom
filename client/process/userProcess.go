package process

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"gocode/chatroom/client/utils"
	"gocode/chatroom/common/message"
	"net"
	"os"
)

type UserProcess struct {
	// 字段
}

// 关联一个用户登录的方法
func (up *UserProcess) Login(userId int, userPwd string) (err error) {
	// 下一个就要开始定协议
	// fmt.Printf("userId = %d userPwd = %s \n", userId, userPwd)

	// return nil

	// 1.链接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	// 延时关闭
	defer conn.Close()
	// 2.准备通过conn发送信息给服务
	var mes message.Message
	mes.Type = message.LoginMesType
	// fmt.Println(mes)
	// 3.创建一个LoginMes 结构体
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd

	// 4.将loginMes序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err= ", err)
		return
	}
	// 5.把data赋给 mes.Data字段
	mes.Data = string(data)
	// 6.将 mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshale err=", err)
		return
	}
	// 7.到这个时候data就是我们要发送的消息
	// 7.1先把 data 的长度发送给服务器
	// 先获取到 data的长度转成一个表示长度的byte切片
	var pkgLen uint32 = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:4], pkgLen)
	// 发送长度
	_, err = conn.Write(buf[:4])
	if err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}
	fmt.Printf("客户端，发送的消息长度=%d 内容=%s\n", len(data), string(data))
	// 发送data
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}
	//休眠20
	// time.Sleep(20 * time.Second)
	// fmt.Println("休眠了20..")
	// 这里还需要处理服务器端返回的消息.
	tf := &utils.Transfer{Conn: conn}
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg(conn) err", err)
		return
	}
	// 将mes的Data反序列化成LoginResMes
	var LoginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &LoginResMes)
	if LoginResMes.Code == 200 {
		// 初始化CurUser
		CurUser.Conn = conn
		CurUser.UserId = userId
		CurUser.UserStatus = message.UserOnline
		// fmt.Println("登录成功")
		// 显示当前在线用户列表，便利LoginResMes.UsersId
		fmt.Println("当前在线用户列表如下：")
		for _, v := range LoginResMes.UsersId {
			if v == userId {
				continue
			}
			fmt.Println("用户id:\t", v)
			// 完成 客户端 onlineUsers初始化
			// 创建User
			user := &message.User{
				UserId: v,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user
		}
		fmt.Print("\n\n")
		// 这里我们还需在客户端启动一个协程
		// 该携程保持和服务端的通信,如果服务器有数据推送给客户端
		// 则接受并显示在客户端的终端
		go serverProcessMes(conn)
		// 1.显示我们登录成功的菜单[循环显示]

		for {
			ShowMenu()
		}
	} else {
		fmt.Println(LoginResMes.Error)
	}
	return
}

// 关联一个用户注册方法
func (up *UserProcess) Register(userId int, userPwd, userName string) (err error) {
	// 下一个就要开始定协议
	// fmt.Printf("userId = %d userPwd = %s \n", userId, userPwd)

	// return nil

	// 1.链接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	// 延时关闭
	defer conn.Close()
	// 2.准备通过conn发送信息给服务
	var mes message.Message
	mes.Type = message.RegisterMesType
	// fmt.Println(mes)
	// 3.创建一个registerMes 结构体
	var registerMes message.RegisterMes
	registerMes.User.UserId = userId
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName

	// 4.将loginMes序列化
	data, err := json.Marshal(registerMes)
	if err != nil {
		fmt.Println("json.Marshal err= ", err)
		return
	}
	// 5.把data赋给 mes.Data字段
	mes.Data = string(data)
	// 6.将 mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshale err=", err)
		return
	}
	// 创建一个Transfer 实例
	tf := &utils.Transfer{Conn: conn}
	// 发送data给服务器端
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("注册发送信息错误err=", err)
	}

	mes, err = tf.ReadPkg() //mes 就是RegisterResMes
	if err != nil {
		fmt.Println("readPkg(conn) err", err)
		return
	}
	// 将mes的Data反序列化成RegisterResMes
	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if registerResMes.Code == 200 {
		fmt.Println("注册成功，请重新登录")
		os.Exit(0)

	} else {
		fmt.Println(registerResMes.Error)
		os.Exit(0)

	}
	return
}
