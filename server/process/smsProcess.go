package processes

import (
	"encoding/json"
	"fmt"
	"net"
	"server/message"
	"server/utils"
)

type SmsProcess struct {
}

func (t *SmsProcess) SendGroupMes(mes *message.Message) {
	//遍历服务器端的onlineUsers map[int]*UserProcess
	//将消息转发出去
	//取出mes的内容 SmsMes
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	for id, up := range userMgr.onlineUsers {
		if id == smsMes.UserId {
			continue
		}
		t.SendMesToEachOnlineUser(data, up.Conn)
	}
}

func (t *SmsProcess) SendMesToEachOnlineUser(data []byte, conn net.Conn) {
	//创建Transfer实例
	tf := &utils.Transfer{
		Conn: conn,
	}
	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("转发消息失败 err=", err)
	}
}
