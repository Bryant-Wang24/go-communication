package main

import (
	"client/process"
	"fmt"
)

// 定义两个变量,一个表示用户id，一个表示用户密码
var userId int
var userPwd string

func main() {
	// 接受用户的选择
	var key int
	for {
		fmt.Println("---------------欢迎登陆多人聊天室---------------")
		fmt.Println("\t\t\t 1 登陆聊天室")
		fmt.Println("\t\t\t 2 注册用户")
		fmt.Println("\t\t\t 3 退出系统")
		fmt.Println("\t\t\t 3 请选择1-3)")
		_, err := fmt.Scanf("%d\n", &key)
		if err != nil {
			return
		}
		switch key {
		case 1:
			fmt.Println("登陆聊天室")
			fmt.Println("请输入用户id")
			_, err := fmt.Scanf("%d\n", &userId)
			if err != nil {
				return
			}
			fmt.Println("请输入用户密码")
			_, err = fmt.Scanf("%s\n", &userPwd)
			if err != nil {
				return
			}
			up := &process.UserProcess{}
			err = up.Login(userId, userPwd)
			if err != nil {
				return
			}
		case 2:
			fmt.Println("注册用户")
		case 3:
			fmt.Println("退出系统")
		default:
			fmt.Println("你的输入有误,请重新输入")
		}
		return
	}
}
