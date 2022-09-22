package process

import (
	"client/message"
	"client/utils"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type UserProcess struct {
}

// Register 用户注册
func (t *UserProcess) Register(userId int, userPwd string, userName string) (err error) {
	fmt.Printf("userId=%d\nuserPwd=%s\nuserName=%s\n", userId, userPwd, userName)
	//1、连接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	// 2、defer关闭conn
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)
	// 3、准备通过conn发送消息给服务器
	var mes message.Message
	mes.Type = message.RegisterMesType
	// 4、创建一个registerMes结构体
	var registerMes message.RegisterMes
	registerMes.User.UserId = userId
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName
	// 5、将registerMes序列化
	data, err := json.Marshal(registerMes)
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
	//	8、到这个时候data就是我们要发送的消息
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
	fmt.Printf("长度=%d 内容=%s\n", len(data), string(data))
	//发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data) fail", err)
		return
	}
	//处理服务端返回的消息
	tf := &utils.Transfer{
		Conn: conn,
	}
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("utils.ReadPkg fail err", err)
		return
	}
	//将mes的data部分反序列化成LoginResMes
	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if registerResMes.Code == 200 {
		//这里还需要在客户端启动一个协程用来保持和服务端的通讯
		//如果服务器有数据推送给客户端，则接受并显示在客户端的终端
		//go serverProcessMes(conn)
		fmt.Println("注册成功")
		os.Exit(0)
		//显示我们的登陆成功的菜单【循环显示】
		//for {
		//	ShowMenu()
		//}
	} else {
		fmt.Println(registerResMes.Error)
		os.Exit(0)
	}
	return
}

// Login 用户登录
func (t *UserProcess) Login(userId int, userPwd string) (err error) {
	fmt.Printf("userId=%d\nuserPwd=%s\n", userId, userPwd)
	// 1、连接到服务器
	conn, err := net.Dial("tcp", "localhost:8889")
	if err != nil {
		fmt.Println("net.Dial err=", err)
		return
	}
	// 2、defer关闭conn
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)
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
	//	8、到这个时候data就是我们要发送的消息
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
	fmt.Printf("长度=%d 内容=%s\n", len(data), string(data))
	//发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data) fail", err)
		return
	}

	//处理服务端返回的消息
	tf := &utils.Transfer{
		Conn: conn,
	}
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("utils.ReadPkg fail err", err)
		return
	}
	//将mes的data部分反序列化成LoginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {
		//登录成功，显示当前在线用户列表
		fmt.Println("当前在线用户列表如下：")
		for _, v := range loginResMes.UsersId {
			//如果要求不显示自己在线，可以增加一个用户id的过滤
			if v == userId {
				continue
			}
			fmt.Println("在线用户id:", v)
			//完成客户端的onlineUsers 初始化
			user := &message.User{
				UserId:     v,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user
		}
		//这里还需要在客户端启动一个协程用来保持和服务端的通讯
		//如果服务器有数据推送给客户端，则接受并显示在客户端的终端
		go serverProcessMes(conn)

		//显示我们的登陆成功的菜单【循环显示】
		for {
			ShowMenu()
		}
	} else {
		fmt.Println(loginResMes.Error)
	}
	return
}
