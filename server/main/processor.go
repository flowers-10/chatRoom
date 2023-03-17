package main

import (
	"fmt"
	"gocode/chatroom/common/message"
	process2 "gocode/chatroom/server/process"
	"gocode/chatroom/server/utils"
	"io"
	"net"
)

// 先创建一个processor的结构体
type Processor struct {
	Conn net.Conn
}

// 编写一个serverProcessMes，根据客户端发送消息种类不同，决定调用哪个函数来处理
func (p *Processor) serverProcessMes(mes *message.Message) (err error) {
	// 测试是否接受客户端发送的群发消息
	fmt.Println("mes=", mes)
	switch mes.Type {
	case message.LoginMesType:
		// 处理登录逻辑
		// 创建一个UserProcess实例
		up := &process2.UserProcess{Conn: p.Conn}
		err = up.ServerProcessLogin(mes)
	case message.RegisterMesType:
		// 处理注册
		up := &process2.UserProcess{Conn: p.Conn}
		err = up.ServerProcessRegister(mes)
	case message.SmsMesType:
		// 处理群聊转发
		// 创建一个SmsProcess 实例完成转发群聊消息
		smsProcess := &process2.SmsProcess{}
		smsProcess.SendGroupMes(mes)
	default:
		fmt.Println("消息类型不存在，无法处理")
	}
	return
}

func (p *Processor) process2() (err error) {
	for {

		// 这里我们将读取数据包，直接封装成一个函数 readPkg(),返回Message，Err
		tf := &utils.Transfer{
			Conn: p.Conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出，服务端也退出。。。")
				return err
			} else {
				fmt.Println("readPkg(conn) err =", err)
				return err
			}
		}
		fmt.Println("mes = ", mes)
		err = p.serverProcessMes(&mes)
		if err != nil {
			return err
		}
	}
}
