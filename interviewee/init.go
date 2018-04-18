package interviewee

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // _ "github.com/lib/pq"
)

type NonOps struct {
	Fid               int    `json:"id"`
	FfullName         string `json:"fullName"`
	FnickName         string `json:"nickName"`
	FphoneNumber      string `json:"phoneNumber"`
	Femail            string `json:"email"`
	Fschool           string `json:"school"`
	Fmajor            string `json:"major"`
	FGPA              string `json:"GPA"`
	Fpurpose          string `json:"purpose"`
	Fcontactpersonid  string `json:"meet"`
	Fposition         string `json:"position"`
	FscheduleTime     string `json:"time"`
	Fjobinfo          string `json:"infoJob"`
	Facquaintance     string `json:"acquaintance"`
	FacquaintanceName string `json:"acquaintanceName"`
	Frelationship     string `json:"relationship"`
	FreferralName     string `json:"referralName"`
	Fstatus           bool   `json:"status"`
	Ftimestamp        string `json:"timestamp"`
}

// type NonOpsG struct {
// 	id               int
// 	FullName         string
// 	NickName         string
// 	PhoneNumber      string
// 	Email            string
// 	School           string
// 	Major            string
// 	GPA              string
// 	Purpose          string
// 	Contactpersonid  string
// 	Positionapply    string
// 	Scheduletime     string
// 	Jobinfo          string
// 	Acquintances     string
// 	Acquintancesname string
// 	Relationship     string
// 	ReferralName     string
// 	Status           bool
// 	Timestamp        string
// }

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
		fmt.Println("Error Binding JSON")
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
		nonops.Fcontactpersonid, nonops.Fposition, nonops.Fjobinfo,
		nonops.Facquaintance, nonops.FscheduleTime, nonops.FacquaintanceName,
		nonops.Frelationship, nonops.FreferralName)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("hello %s, you is %s, your email is %s \n", nonops.FfullName, nonops.FnickName, nonops.Femail)

	c.Writer.Write([]byte(nonops.FacquaintanceName))

}
func ReadData(c *gin.Context) {
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

	rows, err := db.Query(`SELECT * FROM nonops`)
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()
	// println(rows)

	var response []NonOps
	for rows.Next() {
		var nonops NonOps
		if err := rows.Scan(&nonops.Fid, &nonops.FfullName, &nonops.FnickName,
			&nonops.Femail, &nonops.FphoneNumber, &nonops.Fschool, &nonops.Fmajor, &nonops.FGPA,
			&nonops.Fpurpose, &nonops.Fcontactpersonid, &nonops.Fposition,
			&nonops.Fjobinfo, &nonops.Facquaintance, &nonops.FscheduleTime,
			&nonops.FacquaintanceName, &nonops.Frelationship, &nonops.FreferralName, &nonops.Fstatus,
			&nonops.Ftimestamp); err != nil {
			log.Fatal(err)
		}
		response = append(response, nonops)
	}

	b, _ := json.MarshalIndent(response, "", "  ")
	println(string(b))
	c.Writer.Write(b)
}

// if err := rows.Scan(&nonops.id, &nonops.fullName, &nonops.nickName, &nonops.phoneNumber,
// 					&nonops.email, &nonops.school, &nonops.major, &nonops.GPA,
// 					&nonops.purpose, &nonops.contactpersonid, &nonops.positionapply,
// 					&nonops.scheduletime, &nonops.jobinfo, &nonops.acquintances,
// 					&nonops.acquintancesname, &nonops.relationship, &nonops.referralName,
// 					&nonops.timestamp ); err != nil {
// 	log.Fatal(err)
// }
// fmt.Printf("[\"%d\", \"%s\", \"%s\", \"%s\", \"%s\", \"%s\", \"%s\",
// 			 \"%s\", \"%s\", \"%s\", \"%s\", \"%s\", \"%s\", \"%s\",
// 			  \"%s\", \"%s\", \"%s\", \"%s\"],",
// 			  id, fullName, nickName, phoneNumber, email, school, major,
// 			  GPA, purpose, contactpersonid, positionapply, scheduletime,
// 			  jobinfo, acquintances, acquintancesname, relationship, referralName,
// 			  timestamp
// 			)

//fmt.Printf("hello %s, you is %s, your email is %s \n", nonops.FfullName, nonops.FnickName, nonops.Femail)

// c.Writer.Write([]byte(nonops.FacquaintanceName))
