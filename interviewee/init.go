package interviewee

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
	//_ "github.com/lib/pq"
)

type NonOps struct {
	Fid               int
	FfullName         string `json:"fullName"`
	FnickName         string `json:"nickName"`
	FphoneNumber      string `json:"phoneNumber"`
	Femail            string `json:"email"`
	Fschool           string `json:"school"`
	Fmajor            string `json:"major"`
	FGPA              string `json:"GPA"`
	Fpurpose          string `json:"purpose"`
	Fmeet             string `json:"meet"`
	Fposition         string `json:"position"`
	Ftime             string `json:"time"`
	FinfoJob          string `json:"infoJob"`
	Facquaintance     string `json:"acquaintance"`
	FacquaintanceName string `json:"acquaintanceName"`
	Frelationship     string `json:"relationship"`
	FreferralName     string `json:"referralName"`
	Fstatus           bool
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "newpassword"
	dbname   = "hris"
)

func WriteData(c *gin.Context) {
	// fmt.Printf([]byte(r.Body))

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	// var data NonOps
	// c.BindJSON(&data)
	// fmt.Println(data)

	// w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	// w.Header().Set("Access-Control-Allow-Credentials", "true")

	// fmt.Printf("%s", bodyBuffer)
	//json.NewEncoder(w).Encode(bodyBuffer)
	var nonops NonOps
	err := c.BindJSON(&nonops)
	if err != nil {
		fmt.Println("Error JSON Unmarshal")
		fmt.Println(err.Error())
	}

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
				VALUES ($1,$2, $3,$4,
					$5, $6, $7, $8,
					$9, $10, $11,$12,
					$13, $14, $15, $16, true
					)`
	_, err = db.Exec(sqlStatement, nonops.FfullName, nonops.FnickName, nonops.Femail,
		nonops.FphoneNumber, nonops.Fschool, nonops.Fmajor, nonops.FGPA, nonops.Fpurpose,
		nonops.Fmeet, nonops.Fposition, nonops.FinfoJob,
		nonops.Facquaintance, nonops.Ftime, nonops.FacquaintanceName,
		nonops.Frelationship, nonops.FreferralName)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("hello %s, you is %s, your email is %s \n", nonops.FfullName, nonops.FnickName, nonops.Femail)

	c.Writer.Write([]byte(nonops.FacquaintanceName))

}

// func WriteData() {
// 	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
// 		"password=%s dbname=%s sslmode=disable",
// 		host, port, user, password, dbname)
// 	db, err := sql.Open("postgres", psqlInfo)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer db.Close()

// 	err = db.Ping()
// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("Successfully connected!")

// 	sqlStatement := `
// 		INSERT INTO NonOps (fullname, nickname,
// 			email, phone, school, major, gpa, purpose, contactPersonId,
// 			positionApply, jobInfo, acquintances, scheduleTime,
// 			acquintancesName, relationship, referralName, status)
// 		VALUES ("aryl man","aryl", "arylman@gmail.com","081212121212",
// 			"BINUS", "IT", "3.18", "Interview with User",
// 			"Ms.Clarissa", "Software Engineering", "Job Fair","YES",
// 			"10:30PM", "Lala", "Brother/Sister", "Mitha", true
// 			)
// 		RETURNING email`

// 	email := ""
// 	err = db.QueryRow(sqlStatement, "aryl man", "aryl", "arylman@gmail.com", "081212121212",
// 		"BINUS", "IT", "3.18", "Interview with User",
// 		"Ms.Clarissa", "Software Engineering", "Job Fair", "YES",
// 		"10:30PM", "Lala", "Brother/Sister", "Mitha", true).Scan(&email)

// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("New record email is:", email)
// }
