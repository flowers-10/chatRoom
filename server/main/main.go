package main

import (
	"fmt"
	"gocode/chatroom/server/model"
	"net"
	"time"
)

// 携程
func process(conn net.Conn) {
	// 延迟关闭
	defer conn.Close()

	// 调用总控创建
	Processor := &Processor{
		Conn: conn,
	}
	err := Processor.process2()
	if err != nil {
		fmt.Println("客户端和服务器通信携程错误=err", err)
		return
	}

}

// 完成对UserDao的初始化任务
func initUserDao() {
	// pool是全局变量在redis中
	// 注意初始化问题，pool要先初始化才能获得，获得后在启动UserDao
	model.MyUserDao = model.NewUserDao(pool)
}

func main() {
	// 服务器启动初始化redis连接池
	initPool("localhost:6379", 16, 0, 300*time.Second)
	initUserDao()
	//提示信息
	fmt.Println("服务器[新的结构]在8889端口监听....")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")

	if err != nil {
		fmt.Println("net.Listen err=", err)
		return
	}

	defer listen.Close()
	//一旦监听成功，就等待客户端来链接服务器
	for {
		fmt.Println("等待客户端来链接服务器.....")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err=", err)
		}

		//一旦链接成功，则启动一个协程和客户端保持通讯。。
		go process(conn)
	}
}
