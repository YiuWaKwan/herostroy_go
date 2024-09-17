package loginsrv

import (
	"hero_story.go_server/biz_server/base"
	"hero_story.go_server/biz_server/mod/user/userdao"
	"hero_story.go_server/biz_server/mod/user/userdata"
	"hero_story.go_server/comm/async_op"
	"time"
)

// LoginByPasswordAsync 根据用户名称和密码进行登录,
// 将返回一个异步的业务结果
func LoginByPasswordAsync(userName string, password string) *base.AsyncBizResult {
	// 要说下面这两种写法有什么不同么?
	// func LoginByPasswordAsync(userName string, password string, callback func(user *userdata.User)) { ... }
	// func LoginByPasswordAsync(userName string, password string) *base.AsyncBizResult { ... }
	// 第一个是回调方式的, 第二个是返回 Future 方式的,
	// 这两个有什么不一样么?
	// 可以参考其他语言中的 async / await 相关知识...
	//
	if len(userName) <= 0 ||
		len(password) <= 0 {
		return nil
	}

	bizResult := &base.AsyncBizResult{}

	async_op.Process(
		async_op.StrToBindId(userName),
		func() {
			// 通过 DAO 获得用户数据
			user := userdao.GetUserByName(userName)

			nowTime := time.Now().UnixMilli()

			if nil == user {
				// 如果用户数据为空，
				// 则新建数据...
				user = &userdata.User{
					UserName:   userName,
					Password:   password,
					CreateTime: nowTime,
					HeroAvatar: "Hero_Hammer",
				}
			}

			// 更新最后登录时间
			user.LastLoginTime = nowTime
			userdao.SaveOrUpdate(user)

			// 将用户添加到字典
			userdata.GetUserGroup().Add(user)

			bizResult.SetReturnedObj(user)
		},
		nil,
	)

	return bizResult
}
