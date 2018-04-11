package interviewee

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	// "github.com/lib/pq"
)

type NonOps struct {
	id     int
	email  string
	status bool
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "newpassword"
	dbname   = "Hris"
)

func WriteData() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	sqlStatement := `
		INSERT INTO NonOps (fullname, nickname,
			email, phone, school, major, gpa, purpose, contactPersonId,
			positionApply, jobInfo, acquintances, scheduleTime,
			acquintancesName, relationship, referralName, status)
		VALUES ("aryl man","aryl", "arylman@gmail.com","081212121212",
			"BINUS", "IT", "3.18", "Interview with User",
			"Ms.Clarissa", "Software Engineering", "Job Fair","YES",
			"10:30PM", "Lala", "Brother/Sister", "Mitha", true
			)
		RETURNING email`

	email := ""
	err = db.QueryRow(sqlStatement, "aryl man", "aryl", "arylman@gmail.com", "081212121212",
		"BINUS", "IT", "3.18", "Interview with User",
		"Ms.Clarissa", "Software Engineering", "Job Fair", "YES",
		"10:30PM", "Lala", "Brother/Sister", "Mitha", true).Scan(&email)

	if err != nil {
		panic(err)
	}
	fmt.Println("New record email is:", email)
}

func WriteJson(w http.ResponseWriter, r *http.Request) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)

	rows, err := db.Query("SELECT id, email, status FROM NonOps ORDER BY timestamp DESC LIMIT 1")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var email string
		var status bool
		err = rows.Scan(&id, &email, &status)
		if err != nil {
			panic(err)
		}

		NonOps := NonOps{id, email, status}
		js, err := json.Marshal(NonOps)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

}
