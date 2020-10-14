package handler

import (
	db "cloudstorage/db"
	"cloudstorage/util"
	"io/ioutil"
	"net/http"
)

const (
	password_salt = "sd1edawwwqd12eef9cbqu"
)

// 	用户注册
func SignUPHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}

	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	if len(username) < 3 || len(password) < 5 {
		w.Write([]byte("invalid parameter"))
		return
	}

	// 密码加盐
	enc_password := util.Shal([]byte(password + password_salt))
	// 写入数据库
	suc := db.UserSignUp(username, enc_password)
	if suc {
		w.Write([]byte("SUCCESS"))
	} else {
		w.Write([]byte("FAILED"))
	}
}
