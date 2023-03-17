package model

import (
	"encoding/json"
	"fmt"
	"gocode/chatroom/common/message"

	"github.com/garyburd/redigo/redis"
)

// 在服务器启动后，就初始化userDao实例
// 做成全局变量，在需要操作redis时就直接使用即可
var (
	MyUserDao *UserDao
)

// 定义一个UserDao 结构体
// 完成对User 结构体的各种操作

type UserDao struct {
	pool *redis.Pool
}

// 使用工厂模式，创建一个UserDao实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	userDao = &UserDao{
		pool: pool,
	}
	return
}

// UserDao的方法
// 1.根据用户ID返回一个User实例+err
func (this *UserDao) getUserById(conn redis.Conn, id int) (user *User, err error) {

	// 通过给定id去redis查询这个用户
	res, err := redis.String(conn.Do("HGET", "users", id))
	if err != nil {
		// 错误
		if err == redis.ErrNil {
			// 表示users哈希中没有找到对应id
			err = ERROR_USER_NOTEXISTS
		}
		return
	}
	user = &User{}
	// 需要把res反序列化成User实例
	err = json.Unmarshal([]byte(res), user)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	return
}

// 完成登录校验
// 1.Login完成对用户校验
// 2.如果用户id和pwd都正确，返回一个user实例
// 3.如果id或pwd错误，则返回对应错误信息
func (this *UserDao) Login(userId int, userPwd string) (user *User, err error) {

	// 先从UserDao 的连接池中取出一个链接
	conn := this.pool.Get()
	defer conn.Close()
	user, err = this.getUserById(conn, userId)
	if err != nil {
		return
	}

	// 这时证明这个用户是获取到了，验证用户输入的密码和redis里的密码是否一致
	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return
	}

	return
}

// 注册

func (this *UserDao) Register(user *message.User) (err error) {

	// 先从UserDao 的连接池中取出一个链接
	conn := this.pool.Get()
	defer conn.Close()
	_, err = this.getUserById(conn, user.UserId)
	if err == nil {
		err = ERROR_USER_EXISTS
		return
	}

	// 这时说明id不在redis中，可以注册
	data, err := json.Marshal(user) //序列化
	if err != nil {
		return
	}
	// 入库
	_, err = conn.Do("HSET", "users", user.UserId, string(data))
	if err != nil {
		fmt.Println("保存注册用户错误 err=", err)
		return
	}
	return
}
