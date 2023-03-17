package process2

import "fmt"

// UserMgr实例有且只有一个
// 很多地方要用定义一个全局变量
var (
	userMgr *UserMgr
)

type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

// 完成userMgr的初始化工作
func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess,1024),
	}
}

// 完成对onlineUsers添加
func (um *UserMgr) AddOnlineUsers(up *UserProcess) {
	um.onlineUsers[up.UserId] = up
}

// 完成对onlineUsers删除
func (um *UserMgr) DeleteOnlineUsers(userId int) {
	delete(um.onlineUsers,userId)
}

// 完成对onlineUsers查询,返回当前在线所有用户
func (um *UserMgr) GetAllOnlineUsers() map[int]*UserProcess {
	return um.onlineUsers
}

// 完成对onlineUsers查询,根据id返回对应值
func (um *UserMgr) GetOnlineUserById(userId int) (up *UserProcess, err error) {
	// 如何从map中取出一个只，带监测方式
	up, ok :=um.onlineUsers[userId]
	if !ok {
		// 说明查找的用户当前不在线
		err = fmt.Errorf("用户%d不存在",userId)
		return
	} 
	return
}