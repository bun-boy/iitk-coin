package main

import (
  "database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Rollno string `json:"rollno"`
	Name   string `json:"name"`
  Password string `json:"password"`
}

func insertUser(x user) error {
	database, _ :=
		sql.Open("sqlite3", "./User.db")

	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS user (rollno TEXT PRIMARY KEY, name TEXT, password TEXT)")
	statement.Exec()

	statement, _ =
		database.Prepare("INSERT INTO user (rollno, name, password) VALUES (?, ?, ?)")
  _, err := statement.Exec(x.Rollno, x.Name, x.Password)
  return err
}

func getPassword(x user) string {
  database, _ :=
		sql.Open("sqlite3", "./User.db")

  row := database.QueryRow(`SELECT password FROM user WHERE rollno = ?`, x.Rollno)

  var passcode string
  row.Scan(&passcode)
  return passcode
}

func issueToken(rollno string) (string, time.Time, error) {
  var err error

	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Fatal("I am a bad setter :/ Sorry for the error")
	}

	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["userRoll"] = rollno
	expTime := time.Now().Add(time.Minute * 15)
	atClaims["exp"] = expTime.Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESSKEY")))

	if err != nil {
		return "", time.Now(), err
	}
	return token, expTime, err
}

func loginHandler(w http.ResponseWriter, r *http.Request)  {
  if r.URL.Path != "/login" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

  switch r.Method {
  case "POST":
    var account user
    err := json.NewDecoder(r.Body).Decode(&account)

    if err != nil {
      fmt.Println(err)
      http.Error(w, err.Error(), http.StatusBadRequest)
      return
    }

    if account.Rollno == "" || account.Password == "" {
      w.WriteHeader(http.StatusBadRequest)
      if account.Rollno == "" {
        w.Write([]byte("Get a roll number, dude"))
      } else {
        w.Write([]byte("You cannot PASS without a WORD"))
      }
    }

    hash := getPassword(account)
    check_err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(account.Password));

    if check_err != nil {
      w.WriteHeader(500)
      w.Write([]byte("Your word's wrong ;_;"))
      return
    }

    token, expirationTime, err := issueToken(account.Rollno)
    if err != nil {
      w.WriteHeader(401)
      w.Write([]byte("Server error :):"))
      return
    }

    http.SetCookie(w, &http.Cookie {
      Name: "token",
      Value: token,
      Expires: expirationTime,
      HttpOnly: true,
    })

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Successful login! Sigh!"))

  default:
    w.WriteHeader(http.StatusBadRequest)
    fmt.Fprintf(w, "Its login and only POST methods are accepted, you idiot")
  }

}

func signupHandler(w http.ResponseWriter, r *http.Request)  {
  if r.URL.Path != "/signup" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

  switch r.Method {
  case "POST":
    var account user
    err := json.NewDecoder(r.Body).Decode(&account)

    if err != nil {
      fmt.Println(err)
      http.Error(w, err.Error(), http.StatusBadRequest)
      return
    }

    if account.Rollno == "" || account.Password == "" {
      w.WriteHeader(http.StatusBadRequest)
       if account.Rollno == "" {
         w.Write([]byte("Get a roll number, dude"))
       } else {
      w.Write([]byte("You cannot PASS without a WORD"))
      }
      return
    }

    hashed_password, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
    if err != nil {
      w.WriteHeader(401)
      w.Write([]byte("Server error :):"))
      return
    }

    account.Password = string(hashed_password)

    insert_err := insertUser(account)

    if insert_err != nil {
      log.Printf("Body read error, %v", insert_err)
      w.WriteHeader(500)
      w.Write([]byte("You are already signed up :/"))
      return
    }

    fmt.Println("Account created successfully")
    fmt.Fprintf(w, "Account created successfully")

  default:
    w.WriteHeader(http.StatusBadRequest)
    fmt.Fprintf(w, "Its signup and only POST methods are accepted, you idiot")
  }

}

func verifyToken(userToken string) (*jwt.Token, error) {
  errEnv := godotenv.Load()
  if errEnv != nil {
		log.Fatal("I am a bad setter :/ Sorry for the error")
	}
  token, err := jwt.Parse(userToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESSKEY")), nil
	})

  if err != nil {
    return nil, err
  }
  return token, nil
}

func extractData(userToken string) (string, error) {
  token, err := verifyToken(userToken)
  if err != nil {
    return " ", err
  }
  claims, ok := token.Claims.(jwt.MapClaims)
  if ok {
    roll := claims["userRoll"].(string)
    return roll, err
  }
  return " ", err
}

func secretpageHandler(w http.ResponseWriter, r *http.Request)  {
  if r.URL.Path != "/secretpage" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

  switch r.Method {
  case "GET":
    c, err := r.Cookie("token")
    if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Authorization failed"))
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		tokenFromUser := c.Value
		userRoll, err := extractData(tokenFromUser)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Access failed"))
			return
		}
		w.Write([]byte("Welcome "+ userRoll))
		return
  default:
    w.WriteHeader(http.StatusBadRequest)
    fmt.Fprintf(w, "Shh... Its a secret page. Its your turn to GET some data")
  }
}

func main() {
  http.HandleFunc("/signup",signupHandler)
  http.HandleFunc("/login",loginHandler)
  http.HandleFunc("/secretpage",secretpageHandler)
  log.Fatal(http.ListenAndServe(":8000", nil))
}
