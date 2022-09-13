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
