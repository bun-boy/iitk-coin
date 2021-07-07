package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bun-boy/iitk-coin/utils"

	_ "github.com/mattn/go-sqlite3"
)

type transferCoin struct {
	Roll_no string  `json:"rollno"`
	Amount  float64 `json:"amount"`
}

func TransferCoinsHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/transfercoin" {
		resp := &serverResponse{
			Message: "404 Page not found",
		}
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
	}
	tokenFromUser := c.Value
	userRollNo, _, _ := utils.ExtractTokenMetadata(tokenFromUser)
	w.Header().Set("Content-Type", "application/json")
	resp := &serverResponse{
		Message: "",
	}
	switch r.Method {
	case "POST":
		var transferData transferCoin
		err := json.NewDecoder(r.Body).Decode(&transferData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		TransferTo := transferData.Roll_no
		transferAmount := transferData.Amount
		if TransferTo == "" {
			w.WriteHeader(401)
			resp.Message = "Please enter a roll number"
			JsonRes, _ := json.Marshal(resp)
			w.Write(JsonRes)
			return
		}
		err, tax := utils.TransferCoinDb(userRollNo, TransferTo, transferAmount) // withdraw from first user and transfer to second
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		resp.Message = "Transaction of " + fmt.Sprintf("%.2f", transferAmount) + " Sucessfull!  Tax = " + fmt.Sprintf("%.2f", tax)
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	default:
		w.WriteHeader(http.StatusBadRequest)
		resp.Message = "Sorry, only POST requests allowed"
		JsonRes, _ := json.Marshal(resp)
		w.Write(JsonRes)
		return
	}
}
