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
	Conn   net.Conn
	UserId int
}

// NotifyOthersOnlineUser 编写通知所有在线用户的方法。userId 要通知其他的在线用户。自己已经上线
func (t *UserProcess) NotifyOthersOnlineUser(userId int) {
	//遍历onlineUsers,然后一个一个的发送NotifyUserStatusMes
	for id, up := range userMgr.onlineUsers {
		if id == userId {
			continue
		}
		//	开始通知
		up.NotifyMeOnline(userId)
	}
}

func (t *UserProcess) NotifyMeOnline(userId int) {
	//组装NotifyUserStatusMes
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType
	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOnline
	//将notifyUserStatusMes序列化
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//将序列化后的notifyUserStatusMes赋值给mes.Data
	mes.Data = string(data)
	//对mes再次序列化，准备发送
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//发送，创建我们的Transfer实例，发送
	tf := &utils.Transfer{
		Conn: t.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyMeOnline err=", err)
		return
	}
}

// ServerProcessRegister 编写一个函数serverProcessRegister函数，处理注册请求
func (t *UserProcess) ServerProcessRegister(mes *message.Message) (err error) {
	//	先从mes中取出mes.Data,并直接反序列化成RegisterMes
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal fail err=", err)
		return
	}
	//	声明一个resMes
	var resMes message.Message
	resMes.Type = message.RegisterResMesType
	//	声明一个RegisterResMes
	var registerResMes message.RegisterResMes
	//	使用model.MyUserDao到redis数据库去验证
	err = model.MyUserDao.Register(&registerMes.User)
	fmt.Println("err=", err)
	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerResMes.Code = 500
			registerResMes.Error = err.Error()
		} else {
			registerResMes.Code = 505
			registerResMes.Error = "服务器内部错误"
		}
	} else {
		registerResMes.Code = 200
		fmt.Println("注册成功")
	}
	//	将registerResMes序列化
	data, err := json.Marshal(registerResMes)
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

	//使用model.MyUserDao到redis数据库去验证
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
		//这里，因为用户已经登陆成功，我们就把该登陆成功的用户放入到userMgr中
		//将登陆成功的用户的userId赋给this
		t.UserId = loginMes.UserId
		userMgr.AddOnlineUser(t)
		//通知其他的在线用户，我已上线
		t.NotifyOthersOnlineUser(loginMes.UserId)
		//将当前在线用户的id放入到loginResMes.UsersId
		//遍历userMgr.onlineUsers
		for id := range userMgr.onlineUsers {
			loginResMes.UsersId = append(loginResMes.UsersId, id)
			fmt.Println("在线用户id=", id)
		}
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
