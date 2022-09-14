package model

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

// UserDao 定义一个UserDao结构体，完成对User结构体的各种操作
type UserDao struct {
	pool *redis.Pool
}

// 根据用户id返回一个User实例和err
func (t *UserDao) getUserById(conn redis.Conn, id int) (user *User, err error) {
	//	通过给定的id去redis查询这个用户
	res, err := redis.String(conn.Do("HGet", "users", id))
	if err != nil {
		if err == redis.ErrNil { //表示在users 哈希中，没有找到对应的id
			err = ERROR_USER_NOTEXISTS
		}
		return
	}
	user = &User{}
	//	这里我们把res反序列化成User实例
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}
	return
}

// Login 完成登陆的校验
// 如果用户的id和pwd都正确，则返回一个user实例
// 如果用户的id获pwd有错误，则返回对应的错误信息
func (t *UserDao) Login(userId int, userPwd string) (user *User, err error) {
	//	先从UserDao的连接池中取出一根连接
	conn := t.pool.Get()
	defer conn.Close()
	user, err = t.getUserById(conn, userId)
	if err != nil {
		return
	}
	//	这时证明这个用户
	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return
	}
}
