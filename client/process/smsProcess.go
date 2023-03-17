package process

import (
	"encoding/json"
	"fmt"
	"gocode/chatroom/client/utils"
	"gocode/chatroom/common/message"
)

type SmsProcess struct {
}

// 发送群聊消息
func (sp *SmsProcess) SendGroupMes(content string) (err error) {
	// 1创建一个Mes
	var mes message.Message
	mes.Type = message.SmsMesType
	// 2创建一个SmsMes实例
	var smsMes message.SmsMes
	smsMes.Content = content               //内容
	smsMes.UserId = CurUser.UserId         //用户id
	smsMes.UserStatus = CurUser.UserStatus //状态
	// 3序列化
	data, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal err=", err.Error())
		return
	}

	mes.Data = string(data)
	// 4再次序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal err=", err.Error())
		return
	}

	// 5发送服务器
	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}
	// 6发送
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("SendGroupMes tf.WritePkg err=", err.Error())
		return
	}
	return
}
