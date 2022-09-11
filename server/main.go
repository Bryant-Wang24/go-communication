package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"server/message"
)

func readPkg(conn net.Conn) (mes message.Message, err error) {
	buf := make([]byte, 8096)
	fmt.Println("读取客户端发送的数据...")
	//conn.Read在conn没有被关闭的情况下，才会阻塞，如果客户端关闭了conn，就不会阻塞
	_, err = conn.Read(buf[:4])
	if err != nil {
		//err = errors.New("read pkg header err")
		return
	}
	//根据bug[:4]转成一个uint32类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(buf[0:4])
	//	根据pkgLen读取消息内容
	n, err := conn.Read(buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		//err = errors.New("read pkg body err")
		return
	}
	//	把pkgLen反序列化成=》message.Message
	err = json.Unmarshal(buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal err", err)
		return
	}
	return
}

func writePkg(conn net.Conn, data []byte) (err error) {
	//发送一个长度给对方
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
	//发送长度
	n, err := conn.Write(buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write fail", err)
		return
	}
	//	发送data本身
	n, err = conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write fail", err)
		return
	}
	return
}

// 编写一个函数serverProcessLogin函数，处理登陆请求
func serverProcessLogin(conn net.Conn, mes *message.Message) (err error) {
	//	先从mes中取出mes.Data,并直接反序列化成LoginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}
	//	声明一个resMes
	var resMes message.Message
	resMes.Type = message.LoginResMesType
	//	声明一个LoginResMes
	var loginResMes message.LoginResMes

	//	如果用户id=100，密码=123456，认为合法，否则不合法
	if loginMes.UserId == 100 && loginMes.UserPwd == "123456" {
		//	合法
		loginResMes.Code = 200
	} else {
		//	不合法
		loginResMes.Code = 500
		loginResMes.Error = "该用户未注册"
	}
	//	将loginResMes序列化
	data, err := json.Marshal(loginResMes)
	if err != nil {
		fmt.Println("json.Marshal fail err=", err)
		return
	}
	//	将data赋值给resMes
	resMes.Data = string(data)
	//	对resMes 进行序列化，准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal fail err=", err)
		return
	}
	//	发送data，封装到writePkg函数
	err = writePkg(conn, data)
	if err != nil {
		return err
	}
	return
}

// 编写一个 serverProcessMes函数
// 功能：根据客户端发送消息种类不同，决定调用哪个函数来处理
func serverProcessMes(conn net.Conn, mes *message.Message) (err error) {
	switch mes.Type {
	case message.LoginMesType:
		//处理登陆
		err = serverProcessLogin(conn, mes)
	//case : message.Register

	//处理注册
	default:
		fmt.Println("消息类型不存在，无法处理...")
	}
	return
}

// 处理和客户端的通讯
func process(conn net.Conn) {
	//延时关闭coon
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)
	//	循环的读取客户端发送的消息
	for {
		mes, err := readPkg(conn)
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出，服务端也退出...")
				return
			} else {
				fmt.Println("readPkg err=", err)
				return
			}
		}
		fmt.Println("mes=", mes)
		err = serverProcessMes(conn, &mes)
		if err != nil {
			return
		}
	}

}

func main() {
	//	提示信息
	fmt.Println("服务器在8889端口监听...")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Println("net.listen err=", err)
		return
	}
	//如果监听成功，就等待客户端连接服务器
	for {
		fmt.Println("等待客户端连接服务器...")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("net.accept err=", err)
		}
		//旦连接成功，则启动一个协程和客户端保持通讯
		go process(conn)
	}
}
