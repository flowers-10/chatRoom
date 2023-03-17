package process2

import (
	"encoding/json"
	"fmt"
	"gocode/chatroom/client/utils"
	"gocode/chatroom/common/message"
	"net"
)

type SmsProcess struct {
	// 暂时不需要字段
}

// 转发消息
func (sp *SmsProcess) SendGroupMes(msg *message.Message) {
	// 便利服务器端map onlineUsers map[int]*UserProcess
	// 将消息转发出去

	// 取出mes的额SmsMes
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(msg.Data), &smsMes)
	if err != nil {
		fmt.Println(" json.Unmarshal err =", err)
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		fmt.Println(" json.Marshal err =", err)
		return
	}

	for id, up := range userMgr.onlineUsers {
		// 过滤掉自己
		if id == smsMes.UserId {
			continue
		}
		sp.SendMesToEachonlineUser(data, up.Conn)

	}
}

func (sp *SmsProcess) SendMesToEachonlineUser(data []byte, conn net.Conn) {
	// 创建一个Transfer 实例：发送data
	tf := &utils.Transfer{
		Conn: conn,
	}
	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("转发消息失败 err=", err)
	}
}
