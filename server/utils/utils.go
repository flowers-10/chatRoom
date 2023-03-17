package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"gocode/chatroom/common/message"
	"net"
)

// 这里将这些方法关联到结构体中
type Transfer struct {
	// 分析它应该有哪些字段？
	Conn net.Conn
	Buf  [8096]byte //传输使用缓冲

}

// 读取包
func (t *Transfer) ReadPkg() (mes message.Message, err error) {
	// buf := make([]byte, 8096)

	fmt.Println("读取客户端发送数据...")
	_, err = t.Conn.Read(t.Buf[:4])
	if err != nil {
		fmt.Println("conn,Read err =", err)
		return
	}
	// 根据buf[:4]转换成一个uint32类型
	var pkgLen uint32 = binary.BigEndian.Uint32(t.Buf[0:4])
	// 根据 pkglen 读取消息内容
	n, err := t.Conn.Read(t.Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn,Read fail err =", err)
		return
	}

	// 吧pkgLen反序列化成Message
	err = json.Unmarshal(t.Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	return
}

func (t *Transfer) WritePkg(data []byte) (err error) {
	// 先发送一个长度给对方
	var pkgLen uint32 = uint32(len(data))
	// var buf [4]byte
	binary.BigEndian.PutUint32(t.Buf[0:4], pkgLen)
	// 发送长度
	n, err := t.Conn.Write(t.Buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write fial", err)
		return
	}
	n, err = t.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write fial", err)
		return
	}
	return
}
