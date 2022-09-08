package main

import (
	"client/message"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

func login(userId int, userPwd string) (err error) {
	fmt.Printf("userId=%d\nuserPwd=%s\n", userId, userPwd)
	// 1、连接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	// 2、defer关闭conn
	defer conn.Close()
	// 3、准备通过conn发送消息给服务器
	var mes message.Message
	mes.Type = message.LoginMesType
	// 4、创建一个LoginMes结构体
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd
	// 5、将loginMes序列化
	data, err := json.Marshal(loginMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	// 6、将data赋给mes.Data字段
	mes.Data = string(data)
	// 7、将mes进行序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//	7、到这个时候data就是我们要发送的消息
	//先把data的长度发送给服务器
	//	先获取到data的长度=》转成一个表示长度的byte切片.
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
	//	发送长度
	n, err := conn.Write(buf)
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) fail", err)
		return
	}
	//发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data) fail", err)
		return
	}
	return
}
