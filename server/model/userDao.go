package model

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
)

// MyUserDao 定义一个全局变量，在需要和redis操作时，直接使用
var (
	MyUserDao *UserDao
)

// UserDao 定义一个UserDao结构体，完成对User结构体的各种操作
type UserDao struct {
	pool *redis.Pool
}

// NewUserDao 使用工厂模式，创建一个UserDao实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		pool: pool,
	}
	return
}

// 根据用户id返回一个User实例和err
func (t *UserDao) getUserById(conn redis.Conn, id int) (user *User, err error) {
	//	通过给定的id去redis查询这个用户
	res, err := redis.String(conn.Do("HGet", "users", id))
	fmt.Println("redis查询用户res=", res, "redis查询用户err=", err)
	if err != nil {
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

// Register 完成用户注册的校验
// 如果用户id已经存在，则返回对应的错误信息
// 如果用户id不存在，则返回一个user实例
func (t *UserDao) Register(user *User) (err error) {
	//	先从UserDao的连接池中取出一根连接
	conn := t.pool.Get()
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)
	searchUser, err := t.getUserById(conn, user.UserId)
	fmt.Println("user=", user, "err=", err)
	if err != nil {
		if err == redis.ErrNil { //表示在users 哈希中，没有找到对应的id，可以完成注册
			err = nil
			//	将users序列化
			data, err := json.Marshal(user)
			if err != nil {
				return err
			}
			//	将序列化后的user保存到redis
			_, err = conn.Do("HSet", "users", user.UserId, string(data))
			if err != nil {
				fmt.Println("保存注册用户错误 err=", err)
				return err
			}
		}
		return
	}
	//	判断用户是否存在
	if searchUser.UserId == user.UserId {
		err = ERROR_USER_EXISTS
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
	defer func(conn redis.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)
	user, err = t.getUserById(conn, userId)
	//fmt.Println("user=", user, "err=", err)
	if err != nil {
		if err == redis.ErrNil { //表示在users 哈希中，没有找到对应的id,不可登录
			err = ERROR_USER_NOTEXISTS
		}
		return
	}
	//	这时证明这个用户
	if user.UserPwd != userPwd {
		fmt.Println("用户密码错误")
		err = ERROR_USER_PWD
		return
	}
	return
}
