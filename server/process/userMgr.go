package processes

import "fmt"

// userMgr 实例在服务器端有且只有一个，在很多地方都会使用到
// 因此定义为全局变量
var (
	userMgr *UserMgr
)

type UserMgr struct {
	onlineUsers map[int]*UserProcess
}

// 完成对userMgr初始化工作
func init() {
	userMgr = &UserMgr{
		onlineUsers: make(map[int]*UserProcess, 1024),
	}
}

// AddOnlineUser 完成对onlineUsers添加
func (t *UserMgr) AddOnlineUser(up *UserProcess) {
	t.onlineUsers[up.UserId] = up
}

// DeleteOnlineUser 删除
func (t *UserMgr) DeleteOnlineUser(userId int) {
	delete(t.onlineUsers, userId)
}

// GetAllOnlineUser 返回所有当前在线的用户
func (t *UserMgr) GetAllOnlineUser() map[int]*UserProcess {
	return t.onlineUsers
}

// GetOnlineUserById 根据id返回对应的值
func (t *UserMgr) GetOnlineUserById(userId int) (up *UserProcess, err error) {
	//从map取出一个值，带检测方式
	up, ok := t.onlineUsers[userId]
	if !ok { //说明当前查找的用户不在线
		err = fmt.Errorf("用户%d 不存在", userId)
		return
	}
	return
}
