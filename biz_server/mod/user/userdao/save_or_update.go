package userdao

import (
	"hero_story.go_server/biz_server/base"
	"hero_story.go_server/biz_server/mod/user/userdata"
	"hero_story.go_server/comm/log"
)

const sqlSaveOrUpdate = `
insert into t_user ( 
	user_name, password, hero_avatar, curr_hp, create_time, last_login_time
) value (
	?, ?, ?, ?, ?, ?
)
on duplicate key update curr_hp = values(curr_hp), last_login_time = values(Last_login_time)
`

// SaveOrUpdate 保存或者更新用户数据
func SaveOrUpdate(user *userdata.User) {
	if nil == user {
		return
	}

	stmt, err := base.MysqlDB.Prepare(sqlSaveOrUpdate)

	if nil != err {
		log.Error("%+v", err)
		return
	}

	result, err := stmt.Exec(
		user.UserName,
		user.Password,
		user.HeroAvatar,
		user.CurrHp,
		user.CreateTime,
		user.LastLoginTime,
	)

	if nil != err {
		log.Error("%+v", err)
		return
	}

	rowId, err := result.LastInsertId()

	if nil != err {
		log.Error("%+v", err)
		return
	}

	user.UserId = rowId
}
