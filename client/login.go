package main

import "fmt"

func login(userId int, userPwd string) (err error) {

	fmt.Printf("userId=%d\nuserPwd=%s\n", userId, userPwd)
	return nil
}
