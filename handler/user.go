package handler

import (
	dblayer "filestore-server/db"
	"filestore-server/util"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	pwd_salt = "*#890"
)

// SignupHandler: handle user signup request
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := os.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}

	r.ParseForm()
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")

	if len(username) < 3 || len(passwd) < 5 {
		w.Write([]byte("Invalid parameter"))
		return
	}

	enc_passwd := util.Sha1([]byte(passwd + pwd_salt))

	suc := dblayer.UserSignup(username, enc_passwd)
	if suc {
		w.Write([]byte("SUCCESS"))
	} else {
		w.Write([]byte("FAILED"))
	}
	
}

// SignInHandler: handle user signin request
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := os.ReadFile("./static/view/signin.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)

		//http.Redirect(w, r, "/static/view/signin.html", http.StatusFound)
		//return

	}
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	encPasswd := util.Sha1([]byte(password + pwd_salt))

	//1. check the username and password
	pwdChecked := dblayer.UserSignin(username, encPasswd)
	if !pwdChecked{
		w.Write([]byte("FAILED"))
		return
	}
	//2. generate token
	token := GenToken(username)
	upRes := dblayer.UpdateToken(username, token)
	if !upRes {
		w.Write([]byte("FAILED"))
		return
	}

	//3 redirect to home page
	w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
	// resp := util.RespMsg{
	// 	Code: 0,
	// 	Msg:  "OK",
	// 	Data: struct {
	// 		Location string
	// 		Username string
	// 		Token    string
	// 	}{
	// 		Location: "http://" + r.Host + "/static/view/home.html",
	// 		Username: username,
	// 		Token:    token,
	// 	},
	// }
	// w.Write(resp.JSONBytes())



}

func GenToken(username string) string {
	// 40位字符:md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}