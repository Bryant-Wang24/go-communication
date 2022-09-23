package main

import (
	"fmt"
	"io"
	"net"
	"server/message"
	processes "server/process"
	"server/utils"
)

type Processor struct {
	Conn net.Conn
}

// ServerProcessMes 编写一个 serverProcessMes函数
// 功能：根据客户端发送消息种类不同，决定调用哪个函数来处理
func (t *Processor) ServerProcessMes(mes *message.Message) (err error) {
	switch mes.Type {
	case message.LoginMesType:
		//处理登陆
		//创建一个UserProcess实例
		up := &processes.UserProcess{
			Conn: t.Conn,
		}
		err = up.ServerProcessLogin(mes)
	case message.RegisterMesType:
		//处理注册
		up := &processes.UserProcess{
			Conn: t.Conn,
		}
		err = up.ServerProcessRegister(mes)
	case message.SmsMesType:
		//处理群发消息
		up := &processes.SmsProcess{}
		up.SendGroupMes(mes)
	default:
		fmt.Println("消息类型不存在，无法处理...")
	}
	return
}

func (t *Processor) process2() (err error) {
	//	循环的读取客户端发送的消息
	for {
		//这里我们将读取数据包，直接封装成一个函数readPkg(),返回Message,Err
		//创建一个Transfer实例完成读包任务
		tf := &utils.Transfer{
			Conn: t.Conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出，服务端也退出...")
				return err
			} else {
				fmt.Println("readPkg err=", err)
				return err
			}
		}
		fmt.Println("mes=", mes)
		err = t.ServerProcessMes(&mes)
		if err != nil {
			return err
		}
	}
}
