package common

const (
	LoginMesType    = "LoginMes"
	LoginResMesType = "LoginResMes"
)

type Message struct {
	Type string //消息类型
	Data string
}

// 定义两个消息体，后面有需要在增加

type LoginMes struct {
	UserId   int    //用户id
	UserPwd  string //用户密码
	UserName string //用户名
}
type LoginResMes struct {
	Code  int    //状态码
	Error string //错误信息
}
