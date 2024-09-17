package user_lso

import (
	"fmt"
	"hero_story.go_server/biz_server/mod/user/userdao"
	"hero_story.go_server/biz_server/mod/user/userdata"
	"hero_story.go_server/comm/async_op"
)

type UserLso struct {
	*userdata.User
}

func (lso *UserLso) GetLsoId() string {
	return fmt.Sprintf("UserLso_%d", lso.UserId)
}

func (lso *UserLso) SaveOrUpdate() {
	// 通过异步方式保存数据
	async_op.Process(
		int(lso.UserId),
		func() {
			userdao.SaveOrUpdate(lso.User)
		},
		nil,
	)
}
