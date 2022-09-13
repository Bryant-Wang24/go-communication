package model

// User 定义一个用户的结构体
type User struct {
	UserId   int    `json:"userId"`
	UserPwd  string `json:"usePwd"`
	UserName string `json:"userName"`
}
