package db

import (
	"cloudstorage/db/mysql"
	"fmt"
)

// 用户注册
func UserSignUp(username string, password string) bool {
	stmt, err := mysql.DB.Prepare("insert ignore into tbl_user (`user_name`, `user_pwd`) values(?,?)")
	if err != nil {
		fmt.Println("failed to insert, err:" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, password)
	if err != nil {
		fmt.Println("Failed to insert err:" + err.Error())
		return false
	}
	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	// 用户名已经被注册
	return false
}
