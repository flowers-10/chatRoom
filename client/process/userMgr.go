package process

import (
	"fmt"
	"gocode/chatroom/client/model"
	"gocode/chatroom/common/message"
)

// 客户端要维护的map
var onlineUsers map[int]*message.User = make(map[int]*message.User, 10)
// 用户登录成功后完成CurUser初始化
var CurUser model.CurUser


// 在客户端显示在线的用户
func outPutOnlineUser() {
	// 遍历 onlineUsers
	fmt.Println("当前用户在线列表")

	for id, _ := range onlineUsers {
		// 如果不显示自己

		fmt.Println("用户id:\n", id)
	}
}

// 处理返回的NotifyUserStatusMes
func UpdateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {
	// 适当优化
	user, ok := onlineUsers[notifyUserStatusMes.UserId]
	if !ok {
		// 原来没有
		user = &message.User{
			UserId: notifyUserStatusMes.UserId,
		}
	}

	user.UserStatus = notifyUserStatusMes.Status

	onlineUsers[notifyUserStatusMes.UserId] = user

	outPutOnlineUser()
}
