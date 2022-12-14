package main

import (
	"fmt"
	"net"
	"server/model"
	"time"
)

// 处理和客户端的通讯
func process(conn net.Conn) {
	//延时关闭coon
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)
	//调用总控
	processor := &Processor{
		Conn: conn,
	}
	err := processor.process2()
	if err != nil {
		fmt.Println("客户端和服务器通讯协程错误,err=", err)
		return
	}
}

// init 会在main函数之前调用
func init() {
	//当服务启动时，初始化我们的redis连接池
	initPool("localhost:6379", 16, 0, 300*time.Second)
	initUserDao()
}

// 完成对UserDao初始化任务
func initUserDao() {
	//	pool是一个全局的变量
	//	初始化顺序，先initPool，在initUserDao
	model.MyUserDao = model.NewUserDao(pool)
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
