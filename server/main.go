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

// 处理和客户端的通讯
func process(conn net.Conn) {
	//延时关闭coon
	defer conn.Close()
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
