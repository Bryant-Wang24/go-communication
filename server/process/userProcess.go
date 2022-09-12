package processes

import (
	"encoding/json"
	"fmt"
	"net"
	"server/message"
	"server/utils"
)

type UserProcess struct {
	Conn net.Conn
}

// ServerProcessLogin 编写一个函数serverProcessLogin函数，处理登陆请求
func (t *UserProcess) ServerProcessLogin(mes *message.Message) (err error) {
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
	//因为使用分层模式（MVC），先创建一个Transfer实例，然后读取
	tf := &utils.Transfer{
		Conn: t.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		return err
	}
	return
}
