package process

import (
	"encoding/json"
	"fmt"
	"gocode/chatroom/common/message"
)

func OutPutGroupMes(mes *message.Message) {
	// 显示即可
	// 1反序列化
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal err =", err.Error())
		return
	}

	// 显示信息
	info := fmt.Sprintf("用户id:\t %d对大家说：\t %s", smsMes.UserId, smsMes.Content)
	fmt.Println(info)
	fmt.Println()
}
