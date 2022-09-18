package processes

import (
	"encoding/json"
	"fmt"
	"net"
	"server/message"
	"server/model"
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

	////使用model.MyUserDao到redis数据库去验证
	user, err := model.MyUserDao.Login(loginMes.UserId, loginMes.UserPwd)
	if err != nil {
		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误"
		}
	} else {
		loginResMes.Code = 200
		fmt.Println(user, "登陆成功")
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
