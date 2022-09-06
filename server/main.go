package main

import (
	"fmt"
	"net"
)

// 处理和客户端的通讯
func process(conn net.Conn) {

}

func main() {
	//	提示信息
	fmt.Println("服务器在8889端口坚挺...")
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
