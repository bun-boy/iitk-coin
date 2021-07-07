package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bun-boy/iitk-coin/utils"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signup" {
		resp := &serverResponse{
			Message: "404 Page not found",
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}
	resp := &serverResponse{
		Message: "",
	}
	switch r.Method {
	case "POST":
		var user User
		w.Header().Set("Content-Type", "application/json")
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		name := user.Name
		rollno := user.Rollno
		accountType := user.Account_type
		password := user.Password
		if rollno == "" || password == "" || accountType == "" {
			w.WriteHeader(http.StatusBadRequest)
			resp.Message = "Roll No, Password or account type  Cannot be empty"
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}
		hashed_password, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			w.WriteHeader(401)
			resp.Message = "Server error"
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
		}
		write_err := utils.WriteUserToDb(name, rollno, string(hashed_password), accountType)
		if write_err != nil {
			w.WriteHeader(500)
			resp.Message = "Roll no already exists"
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}
		w.WriteHeader(http.StatusOK)
		resp.Message = "Account created successfully! Head to login"
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		resp.Message = "Only POST requests allowed"
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}
}
