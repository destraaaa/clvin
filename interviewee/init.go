package interviewee

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gosample/chart"
	_ "github.com/lib/pq" // _ "github.com/lib/pq"
)

type Candidate struct {
	Fid               int       `json:"id"`
	FfullName         string    `json:"fullName"`
	FnickName         string    `json:"nickName"`
	FphoneNumber      string    `json:"phoneNumber"`
	Femail            string    `json:"email"`
	Fschool           string    `json:"school"`
	Fmajor            string    `json:"major"`
	FGPA              string    `json:"GPA"`
	Fpurpose          string    `json:"purpose"`
	Fcontactpersonid  string    `json:"meet"`
	Fposition         string    `json:"position"`
	FscheduleTime     string    `json:"time"`
	Fjobinfo          string    `json:"infoJob"`
	Facquaintance     string    `json:"acquaintance"`
	FacquaintanceName string    `json:"acquaintanceName"`
	Frelationship     string    `json:"relationship"`
	FreferralName     string    `json:"referralName"`
	Fstatus           bool      `json:"status"`
	FformType         string    `json:"formType"`
	Fprogress         int       `json:"progress,omitempty"`
	Fprog             string    `json:"statProgress,omitempty"`
	Ftimestamp        time.Time `json:"timestamp"`
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
	var nonops Candidate
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
				INSERT INTO Candidate (fullname, nickname,
					email, phone, school, major, gpa, purpose, contactPersonId,
					positionApply, jobInfo, acquaintance, scheduleTime,
					acquaintanceName, relationship, referralName,formType, progress,status)
				VALUES ($1,$2, $3,$4,
					$5, $6, $7, $8,
					$9, $10, $11,$12,
					$13, $14, $15, $16, $17,$18,true
					)`
	_, err = db.Exec(sqlStatement, nonops.FfullName, nonops.FnickName, nonops.Femail,
		nonops.FphoneNumber, nonops.Fschool, nonops.Fmajor, nonops.FGPA, nonops.Fpurpose,
		nonops.Fcontactpersonid, nonops.Fposition, nonops.Fjobinfo,
		nonops.Facquaintance, nonops.FscheduleTime, nonops.FacquaintanceName,
		nonops.Frelationship, nonops.FreferralName, nonops.FformType, nonops.Fprogress)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("hello %s, you is %s, your email is %s \n", nonops.FfullName, nonops.FnickName, nonops.Femail)

	// c.Writer.Write([]byte(nonops.FacquaintanceName))

}
func ReadDataNon(c *gin.Context) {
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

	result := chart.Filtering()

	sqlStatement := `SELECT id, fullname, nickname,
	email, (CASE WHEN progress=1 THEN 'no status'
				WHEN progress=2 THEN 'reject'
				WHEN progress=3 THEN 'approved'
				ELSE 'on progress' END)as progress, phone, school, major, gpa,
	purpose, contactpersonid, positionapply,
	jobinfo, acquaintance, scheduletime,
	acquaintancename, relationship, referralName, status,
	logtimestamps FROM candidate WHERE formtype = 'Non Operational Form' ` + result

	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()
	// println(rows)

	var response []Candidate
	for rows.Next() {
		var nonops Candidate
		if err := rows.Scan(&nonops.Fid, &nonops.FfullName, &nonops.FnickName,
			&nonops.Femail, &nonops.Fprog, &nonops.FphoneNumber, &nonops.Fschool, &nonops.Fmajor, &nonops.FGPA,
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

func ReadRejected(c *gin.Context) {
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

	// result := chart.Filtering()

	sqlStatement := `SELECT id, fullname, nickname,
	email, phone, school, major, gpa,
	purpose, contactpersonid, positionapply,
	jobinfo, acquaintance, scheduletime,
	acquaintancename, relationship, referralName, status,
	logtimestamps FROM candidate WHERE progress = 2`

	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()
	// println(rows)

	var response []Candidate
	for rows.Next() {
		var nonops Candidate
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

func ReadDataOps(c *gin.Context) {
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

	result := chart.Filtering()
	sqlStatement := `SELECT id, fullname, nickname,
	email, (CASE WHEN progress=1 THEN 'no status'
		WHEN progress=2 THEN 'reject'
		WHEN progress=3 THEN 'approved'
		ELSE 'on progress' END)as progress, phone, school, purpose, contactpersonid,
	positionapply, jobinfo, acquaintance, scheduletime,
	acquaintancename, relationship, referralName, status,
	logtimestamps FROM candidate WHERE formtype = 'Operational Form' ` + result

	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()
	// println(rows)

	var response []Candidate
	for rows.Next() {
		var ops Candidate
		if err := rows.Scan(&ops.Fid, &ops.FfullName, &ops.FnickName,
			&ops.Femail, &ops.Fprog, &ops.FphoneNumber, &ops.Fschool, &ops.Fpurpose,
			&ops.Fcontactpersonid, &ops.Fposition, &ops.Fjobinfo,
			&ops.Facquaintance, &ops.FscheduleTime, &ops.FacquaintanceName,
			&ops.Frelationship, &ops.FreferralName, &ops.Fstatus,
			&ops.Ftimestamp); err != nil {
			log.Fatal(err)
		}
		response = append(response, ops)
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
