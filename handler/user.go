package handler

import (
	"net/http"
	"os"
	"filestore-server/util"
	dblayer "filestore-server/db"
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

func SignInHandler(w http.ResponseWriter, r *http.Request) {
		
	//1. check the username and password


	//2. generate token


	//3



}