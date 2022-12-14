package process

import (
	"client/message"
	"client/model"
	"encoding/json"
	"fmt"
)

// 客户端要维护的map
var onlineUsers map[int]*message.User = make(map[int]*message.User, 10)
var CurUser model.CurUser // 在用户登录成功后，完成对CurUser的初始化

// 在客户端显示当前在线用户
func outputOnlineUser() {
	//遍历onlineUsers
	fmt.Println("当前在线用户列表如下:")
	for id := range onlineUsers {
		fmt.Println("在线用户id:", id)
	}
}

// 处理返回的NotifyUserStatusMes
func updateUserStatus(notifyUserStatusMes *message.NotifyUserStatusMes) {
	user, ok := onlineUsers[notifyUserStatusMes.UserId]
	if !ok {
		user = &message.User{
			UserId: notifyUserStatusMes.UserId,
		}
	}
	user.UserStatus = notifyUserStatusMes.Status
	onlineUsers[notifyUserStatusMes.UserId] = user
	outputOnlineUser()
}

// 展示群发消息
func outputGroupMes(mes *message.Message) {
	//反序列化mes.Data
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal err", err)
		return
	}
	//显示群聊消息
	fmt.Printf("用户id:\t%d 对大家说:\t%s\n", smsMes.UserId, smsMes.Content)
}
