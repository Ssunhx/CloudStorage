package db

import (
	"cloudstorage/db/mysql"
	"fmt"
)

type User struct {
	Username     string
	Email        string
	Phone        string
	SignupAt     string
	LastActiveAt string
	Status       int
}

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

func UserSignin(username string, encpwd string) bool {
	stmt, err := mysql.DB.Prepare("select * from tbl_user where user_name=? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if rows == nil {
		return false
	}
	pRows := mysql.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		return true
	}
	return false
}

// 更新用户 token
func UpdateToken(username string, token string) bool {
	stmt, err := mysql.DB.Prepare("replace into tbl_user_token(`user_name`, `user_token`) value (?, ?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

//
func GetUserInfo(username string) (User, error) {
	user := User{}

	stmt, err := mysql.DB.Prepare("select user_name, signup_at from tbl_user where user_name = ?")
	if err != nil {
		fmt.Println(err.Error())
		return user, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil {
		return user, err
	}
	return user, nil
}
