package process

import (
	"client/message"
	"client/utils"
	"encoding/json"
	"fmt"
)

type SmsProcess struct {
}

func (t *SmsProcess) SendGroupMes(content string) (err error) {
	// 1. 创建一个Mes
	var mes message.Message
	mes.Type = message.SmsMesType
	// 2. 创建一个SmsMes实例
	var smsMes message.SmsMes
	smsMes.Content = content
	smsMes.UserId = CurUser.UserId
	smsMes.UserStatus = CurUser.UserStatus
	// 3. 将SmsMes序列化
	data, err := json.Marshal(smsMes)
	if err != nil {
		return
	}
	// 4. 将data赋值给mes.Data字段
	mes.Data = string(data)
	// 5. 对mes再次序列化，准备发送
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal fail err=", err)
		return
	}
	// 6. 发送，创建Transfer实例，发送
	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("发送消息失败 err=", err)
		return
	}
	return
}
