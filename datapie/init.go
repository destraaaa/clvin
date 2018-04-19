package datapie

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // _ "github.com/lib/pq"
)

type pieChart struct {
	Labels string `json:"labels"`
	Series int32  `json:"series"`
	// names
	// types
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "newpassword"
	dbname   = "hris"
)

func ShowPie(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

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

	schooldb, err := db.Query(`SELECT LOWER(school), count(school) FROM nonops GROUP BY LOWER(school)`)
	if err != nil {
		log.Panic(err)
	}
	defer schooldb.Close()
	// println(rows)

	var dataSchool []pieChart

	for schooldb.Next() {
		var school pieChart
		if err := schooldb.Scan(&school.Labels, &school.Series); err != nil {
			log.Fatal(err)
		}
		dataSchool = append(dataSchool, school)
	}

	b, _ := json.MarshalIndent(dataSchool, "", "  ")
	println(string(b))
	c.Writer.Write(b)
}

// var a = []int{1,2,1,2,3,2}

// func main(){
//     m := make(map[string]int)
//     for _, v := range a {
//         m[v]++
//     }

//     fmt.Println(m)
// }
