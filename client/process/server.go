package process

import (
	"encoding/json"
	"fmt"
	"gocode/chatroom/client/utils"
	"gocode/chatroom/common/message"
	"net"
	"os"
)
// 显示登录成功后的界面
func ShowMenu() {
	fmt.Println("---------恭喜xxx登录成功---------")
	fmt.Println("---------1.显示在线用户列表---------")
	fmt.Println("---------2.发送消息---------")
	fmt.Println("---------3.信息列表---------")
	fmt.Println("---------4.退出系统---------")
	fmt.Println("请选择1-4：")
	var key int
	var content string

	// 因为SmsProcess实例需要多次使用，定义在外部
	smsProcess := &SmsProcess{}
	fmt.Scanf("%d\n", &key)
	switch key {
	case 1:
		// fmt.Println("显示在线用户列表")
		outPutOnlineUser()
	case 2:
		fmt.Println("请输入内容")
		fmt.Scanf("%s\n", &content)
		smsProcess.SendGroupMes(content)
	case 3:
		fmt.Println("信息列表")
	case 4:
		fmt.Println("你退出系统。。。。")
		os.Exit(0)
	default:
		fmt.Println("输入选项不正确。。")
	}

}

// 和服务器端保持通讯
func serverProcessMes(Conn net.Conn) {
	// 创建一个transfer实例不停的读取服务器发送的信息
	tf := &utils.Transfer{Conn: Conn}
	for {
		fmt.Println("客户端正在等待读取服务器发送消息")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.ReadPkg() err服务器挂掉了", err)
			return
		}
		// 如果读取到消息,下一步
		switch mes.Type {
		case message.NotifyUserStatusMesType: //有人上线了
			// 1.取出 NotifyUserStatusMes
			var notifyUserStatusMes message.NotifyUserStatusMes
			json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)

			// 2.把这个用的信息状态保存到客户map[int]User中
			UpdateUserStatus(&notifyUserStatusMes)
		// 处理群发消息
		case message.SmsMesType:
			OutPutGroupMes(&mes)
		default:
			fmt.Println("服务器端返回了未知消息类型")
		}
		// fmt.Printf("mes = %v\n", mes)
	}
}
