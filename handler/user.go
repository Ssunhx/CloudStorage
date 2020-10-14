package handler

import (
	db "cloudstorage/db"
	"cloudstorage/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
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

// 登录接口
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signin.html")
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
	enc_password := util.Shal([]byte(password + password_salt))
	// 1、校验用户名密码
	pwdChecked := db.UserSignin(username, enc_password)
	if !pwdChecked {
		w.Write([]byte("FAILED"))
		return
	}
	// 2、生成访问凭证
	token := GenToken(username)
	upres := db.UpdateToken(username, token)
	if !upres {
		w.Write([]byte("FAILED"))
		return
	}
	// 3、登录成功后重定向到首页
	//w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	w.Write(resp.JSONBytes())
}

func GenToken(username string) string {
	// md5(username + timestamp + token_salt )+ tamestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	token_Prefix := util.MD5([]byte(username + "_token_salt"))
	return token_Prefix + ts[:8]
}
