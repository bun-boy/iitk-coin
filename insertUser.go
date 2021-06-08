package main

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type user struct {
	rollno string
	name   string
}

func newUser(roll string, nm string) *user {
	u := user{rollno: roll, name: nm}
	return &u
}

func insertUser(x user) {
	database, _ :=
		sql.Open("sqlite3", "./User.db")
	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, rollno TEXT, name TEXT)")
	statement.Exec()
  row := database.QueryRow("select rollno from people where rollno = ?", x.rollno)
  temp := ""
  row.Scan(&temp)
  if temp != "" {
    return
  }
	statement, _ =
		database.Prepare("INSERT INTO people (rollno, name) VALUES (?, ?)")
	statement.Exec(x.rollno, x.name)
}

func main() {

	p := newUser("190995", "Yash Burnwal")
	insertUser(*p)

	database, _ :=
		sql.Open("sqlite3", "./User.db")

	rows, _ :=
		database.Query("SELECT id, rollno, name FROM people")
	var id int
	var rollno string
	var name string
	for rows.Next() {
		rows.Scan(&id, &rollno, &name)
		fmt.Println(strconv.Itoa(id) + ": " + rollno + " " + name)
	}
}
