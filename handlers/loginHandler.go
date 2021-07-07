package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bun-boy/iitk-coin/utils"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Name         string `json:"name"`
	Rollno       string `json:"rollno"`
	Password     string `json:"password"`
	Account_type string `json:"account_type"`
}

type serverResponse struct {
	Message string `json:"message"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		resp := &serverResponse{
			Message: "404 Page not found",
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	resp := &serverResponse{
		Message: "",
	}
	switch r.Method {
	case "POST":
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		rollno := user.Rollno
		password := user.Password
		hashedPassword := utils.Get_hashed_password(rollno)
		if hashedPassword == "" {
			w.WriteHeader(401)
			resp.Message = "User does not exist"
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
			w.WriteHeader(500)
			resp.Message = "Incorrect password! Please try again"
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}
		_, accountType, _ := utils.GetUserFromRollNo(rollno)
		token, expirationTime, err := utils.CreateToken(rollno, accountType)
		if err != nil {
			w.WriteHeader(401)
			resp.Message = "Server Error"
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			Expires:  expirationTime,
			HttpOnly: true,
		})
		w.WriteHeader(http.StatusOK)
		resp.Message = "Successfully logged in!"
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		resp.Message = "Sorry, only POST requests allowed!"
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}

}
