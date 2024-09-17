package userdao

import (
	"hero_story.go_server/biz_server/base"
	"hero_story.go_server/biz_server/mod/user/userdata"
	"hero_story.go_server/comm/log"
)

const sqlGetUserByName = `select user_id, user_name, password, hero_avatar, curr_hp from t_user where user_name = ?`

// GetUserByName 根据用户名获得用户数据
func GetUserByName(userName string) *userdata.User {
	if len(userName) <= 0 {
		return nil
	}

	row := base.MysqlDB.QueryRow(sqlGetUserByName, userName)

	if nil == row {
		return nil
	}

	user := &userdata.User{}

	err := row.Scan(
		&user.UserId,
		&user.UserName,
		&user.Password,
		&user.HeroAvatar,
		&user.CurrHp,
	)

	if nil != err {
		log.Error("%+v", err)
		return nil
	}

	return user
}
