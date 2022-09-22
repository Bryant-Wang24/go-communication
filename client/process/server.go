package process

import (
	"client/message"
	"client/utils"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

func ShowMenu() {
	fmt.Println("-------恭喜xxx登陆成功---------")
	fmt.Println("-------1.显示在线用户列表---------")
	fmt.Println("-------2.发送消息---------")
	fmt.Println("-------3.消息列表---------")
	fmt.Println("-------4.退出系统---------")
	fmt.Println("请选择1-4")
	var key int
	_, err := fmt.Scanf("%d\n", &key)
	if err != nil {
		return
	}
	switch key {
	case 1:
		//fmt.Println("显示在线用户列表")
		outputOnlineUser()
	case 2:
		fmt.Println("发送消息")
	case 3:
		fmt.Println("消息列表")
	case 4:
		fmt.Println("退出系统")
		os.Exit(0)
	default:
		fmt.Println("你输入的选项不正确")
	}
}

// 和服务器端保持通讯
func serverProcessMes(conn net.Conn) {
	//	创建一个transfer实例，不停的读取服务器发送的消息
	tf := &utils.Transfer{
		Conn: conn,
	}
	for {
		fmt.Println("客户端正在读取服务器发送的消息")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.ReadPkg fail err=", err)
			return
		}
		//	如果读取到消息
		switch mes.Type {
		case message.NotifyUserStatusMesType: //有人上线
			//1、取出.NotifyUserStatusMes
			var notifyUserStatusMes message.NotifyUserStatusMes
			err := json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
			if err != nil {
				return
			}
			updateUserStatus(&notifyUserStatusMes)
			//2.把这个用户的信息，状态保存到客户map[int]User中
		default:
			fmt.Println("服务器返回了消息")
		}

		fmt.Printf("mes=%v", mes)
	}
}
