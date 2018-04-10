package interviewee

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	// "github.com/lib/pq"
)
type NonOps struct{
	fullname string,
	nickname string,
	email string,
	phone string,
	school string,
	major string,
	gpa string,
	purpose string,
	contactperson string,
	positionapply string, 
	jobinfo string,
	acquintances string, 
	scheduletime string,
	acquintancesname string, 
	relationship string, 
	referralname string, 
	status bool
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "newpassword"
	dbname   = "Hris"
)

// type ServerConfig struct {
// 	Name string
// }

// type Config struct {
// 	Server ServerConfig
// }

// type HelloWorldModule struct {
// 	cfg       *Config
// 	something string
// 	stats     *expvar.Int
// }

func RetrieveData() {
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
			email, phone, school, major, gpa, purpose, contactperson,
			positionapply, jobinfo, acquintances, scheduletime,
			acquintancesname, relationship, referralname, status)
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
	fmt.Println("New record ID is:", email)
}

// func (hlm *HelloWorldModule) SayHelloWorld(w http.ResponseWriter, r *http.Request) {
// 	span, ctx := opentracing.StartSpanFromContext(r.Context(), r.URL.Path)
// 	defer span.Finish()

// 	hlm.stats.Add(1)
// 	hlm.someSlowFuncWeWantToTrace(ctx, w)
// }

// func (hlm *HelloWorldModule) someSlowFuncWeWantToTrace(ctx context.Context, w http.ResponseWriter) {
// 	span, ctx := opentracing.StartSpanFromContext(ctx, "someSlowFuncWeWantToTrace")
// 	defer span.Finish()

// 	w.Write([]byte("Hello Kak " + hlm.something))
// }

func foo(w http.ResponseWriter, r *http.Request) {
	profile := Profile{"Alex", []string{"snowboarding", "programming"}}

	js, err := json.Marshal(profile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
