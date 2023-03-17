package main

import (
	"fmt"
	"gocode/chatroom/client/process"
	"os"
)

// 定义量个变量一个表示用户id 一个表示用户密码
var userId int
var userPwd string
var userName string

func main() {
	// 接受用户的选择
	var key int
	// 判断是否连续显示菜单
	// var loop = true

	for true {
		fmt.Println("----------------欢迎登陆多人聊天系统------------")
		fmt.Println("\t\t\t 1 登陆聊天室")
		fmt.Println("\t\t\t 2 注册用户")
		fmt.Println("\t\t\t 3 退出系统")
		fmt.Println("\t\t\t 请选择(1-3):")

		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("登陆聊天室")
			fmt.Println("请输入用户的id")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户的密码")
			fmt.Scanf("%s\n", &userPwd)
			// 完成登录
			// 1.创建一个UserProcess实例
			up := &process.UserProcess{}
			up.Login(userId,userPwd)
			// loop = false
		case 2:
			fmt.Println("注册用户")
			fmt.Println("请输入用户id：")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户的密码")
			fmt.Scanf("%s\n", &userPwd)
			fmt.Println("请输入用户名")
			fmt.Scanf("%s\n", &userName)
			// 2.完成注册请求
			up := &process.UserProcess{}
			up.Register(userId,userPwd,userName)
		case 3:
			fmt.Println("退出系统")
			//loop = false
			os.Exit(0)
		default:
			fmt.Println("你的输入有误，请重新输入")
		}

	}
	//更加用户的输入，显示新的提示信息
	// if key == 1 {
	// 	//说明用户要登陆
	// 	fmt.Println("请输入用户的id")
	// 	fmt.Scanf("%d\n", &userId)
	// 	fmt.Println("请输入用户的密码")
	// 	fmt.Scanf("%s\n", &userPwd)

	// 	// 因为使用了分层，我们创建

	// 	//先把登陆的函数，写到另外一个文件， 比如login.go
	// 	// login(userId, userPwd)

	// } else if key == 2 {
	// 	fmt.Println("进行用户注册的逻辑....")
	// }

}
