package interviewee

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/destraaaa/clvin/chart"
	"github.com/destraaaa/clvin/env"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type Candidate struct {
	Fid               int       `json:"id"`
	FPIC              int       `json:"pic"`
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
	FupdatedDate      time.Time `json:"updatedDate"`
}

func WriteData(c *gin.Context) {
	// fmt.Printf([]byte(r.Body))
	dbconfig := env.GetConfig().Database
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	var nonops Candidate
	err := c.BindJSON(&nonops)
	if err != nil {
		fmt.Println("Error Binding JSON")
		fmt.Println(err.Error())
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbconfig.Host, dbconfig.Port, dbconfig.User, dbconfig.Password, dbconfig.Name)
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

}
func ReadDataNon(c *gin.Context) {
	dbconfig := env.GetConfig().Database
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbconfig.Host, dbconfig.Port, dbconfig.User, dbconfig.Password, dbconfig.Name)
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
	email, (CASE WHEN progress=1 THEN 'NO STATUS' 
		WHEN progress=2 THEN 'REJECT' 
		WHEN progress=3 THEN 'APPROVED' 
		WHEN progress=4 THEN 'ON PROGRESS' 
		WHEN progress=5 THEN 'OFFERING - ACCEPTED' 
		WHEN progress=6 THEN 'OFFERING - DECLINED' 
		WHEN progress=7 THEN 'OFFERING - CANCEL' 
		WHEN progress=8 THEN 'HOLD' 
		WHEN progress=9 THEN 'HOLD - REJECT' 
		ELSE 'CLOSED' END)as progress , phone, school, major, gpa,
	purpose, contactpersonid, positionapply,
	jobinfo, acquaintance, scheduletime,
	acquaintancename, relationship, referralName, status,
	logtimestamps, updatedDate  FROM candidate WHERE formtype = 'Non Operational Form' ` + result

	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()

	var response []Candidate
	var nickname, phone, school, jobinfo, acquaintanceName, relationship, referralName, major, GPA sql.NullString
	var updated pq.NullTime
	for rows.Next() {
		var nonops Candidate
		if err := rows.Scan(&nonops.Fid, &nonops.FfullName, &nickname,
			&nonops.Femail, &nonops.Fprog, &phone, &school, &major, &GPA,
			&nonops.Fpurpose, &nonops.Fcontactpersonid, &nonops.Fposition,
			&jobinfo, &nonops.Facquaintance, &nonops.FscheduleTime,
			&acquaintanceName, &relationship, &referralName, &nonops.Fstatus,
			&nonops.Ftimestamp, &updated); err != nil {
			log.Fatal(err)
		}

		if nickname.Valid {
			temp, _ := nickname.Value()
			strTemp, _ := temp.(string)
			nonops.FnickName = strTemp
		}
		if phone.Valid {
			temp, _ := phone.Value()
			strTemp, _ := temp.(string)
			nonops.FphoneNumber = strTemp
		}
		if school.Valid {
			temp, _ := school.Value()
			strTemp, _ := temp.(string)
			nonops.Fschool = strTemp
		}
		if jobinfo.Valid {
			temp, _ := jobinfo.Value()
			strTemp, _ := temp.(string)
			nonops.Fjobinfo = strTemp
		}
		if acquaintanceName.Valid {
			temp, _ := acquaintanceName.Value()
			strTemp, _ := temp.(string)
			nonops.FacquaintanceName = strTemp
		}
		if relationship.Valid {
			temp, _ := relationship.Value()
			strTemp, _ := temp.(string)
			nonops.Frelationship = strTemp
		}
		if major.Valid {
			temp, _ := major.Value()
			strTemp, _ := temp.(string)
			nonops.Fmajor = strTemp
		}
		if GPA.Valid {
			temp, _ := GPA.Value()
			strTemp, _ := temp.(string)
			nonops.FGPA = strTemp
		}
		if referralName.Valid {
			temp, _ := referralName.Value()
			strTemp, _ := temp.(string)
			nonops.FreferralName = strTemp
		}
		if updated.Valid {
			temp, _ := updated.Value()
			strTemp, _ := temp.(time.Time)
			nonops.FupdatedDate = strTemp
		}
		response = append(response, nonops)
	}

	b, _ := json.MarshalIndent(response, "", "  ")
	println(string(b))
	c.Writer.Write(b)
}

func ReadRejected(c *gin.Context) {
	dbconfig := env.GetConfig().Database
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbconfig.Host, dbconfig.Port, dbconfig.User, dbconfig.Password, dbconfig.Name)
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

	var response []Candidate
	var nickname, phone, school, jobinfo, acquaintanceName, relationship, referralName, major, GPA sql.NullString

	for rows.Next() {
		var nonops Candidate
		if err := rows.Scan(&nonops.Fid, &nonops.FfullName, &nickname,
			&nonops.Femail, &phone, &school, &major, &GPA,
			&nonops.Fpurpose, &nonops.Fcontactpersonid, &nonops.Fposition,
			&jobinfo, &nonops.Facquaintance, &nonops.FscheduleTime,
			&acquaintanceName, &relationship, &referralName, &nonops.Fstatus,
			&nonops.Ftimestamp); err != nil {
			log.Fatal(err)
		}

		if nickname.Valid {
			temp, _ := nickname.Value()
			strTemp, _ := temp.(string)
			nonops.FnickName = strTemp
		}
		if phone.Valid {
			temp, _ := phone.Value()
			strTemp, _ := temp.(string)
			nonops.FphoneNumber = strTemp
		}
		if school.Valid {
			temp, _ := school.Value()
			strTemp, _ := temp.(string)
			nonops.Fschool = strTemp
		}
		if jobinfo.Valid {
			temp, _ := jobinfo.Value()
			strTemp, _ := temp.(string)
			nonops.Fjobinfo = strTemp
		}
		if acquaintanceName.Valid {
			temp, _ := acquaintanceName.Value()
			strTemp, _ := temp.(string)
			nonops.FacquaintanceName = strTemp
		}
		if relationship.Valid {
			temp, _ := relationship.Value()
			strTemp, _ := temp.(string)
			nonops.Frelationship = strTemp
		}
		if major.Valid {
			temp, _ := major.Value()
			strTemp, _ := temp.(string)
			nonops.Fmajor = strTemp
		}
		if GPA.Valid {
			temp, _ := GPA.Value()
			strTemp, _ := temp.(string)
			nonops.FGPA = strTemp
		}
		if referralName.Valid {
			temp, _ := referralName.Value()
			strTemp, _ := temp.(string)
			nonops.FreferralName = strTemp
		}

		response = append(response, nonops)
	}

	b, _ := json.MarshalIndent(response, "", "  ")
	println(string(b))
	c.Writer.Write(b)
}
func UpdateCandidate(c *gin.Context) {
	dbconfig := env.GetConfig().Database
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	var ops Candidate
	err := c.BindJSON(&ops)
	if err != nil {
		fmt.Println("Error Binding JSON")
		fmt.Println(err.Error())
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbconfig.Host, dbconfig.Port, dbconfig.User, dbconfig.Password, dbconfig.Name)
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

	sqlStatement := `UPDATE candidate SET progress = $1, updateddate = $3, pic = $4 WHERE id = $2`
	_, err = db.Exec(sqlStatement, ops.Fprogress, ops.Fid, ops.FupdatedDate, ops.FPIC)
	if err != nil {
		panic(err)
	}
}

func EmailCandidate(c *gin.Context) {
	dbconfig := env.GetConfig().Database
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbconfig.Host, dbconfig.Port, dbconfig.User, dbconfig.Password, dbconfig.Name)
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
	sqlStatement := `SELECT email FROM candidate`

	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()

	var response []Candidate
	for rows.Next() {
		var ops Candidate
		if err := rows.Scan(&ops.Femail); err != nil {
			log.Fatal(err)
		}
		response = append(response, ops)
	}

	b, _ := json.MarshalIndent(response, "", "  ")
	println(string(b))
	c.Writer.Write(b)
}
func ReadDataOps(c *gin.Context) {
	dbconfig := env.GetConfig().Database
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbconfig.Host, dbconfig.Port, dbconfig.User, dbconfig.Password, dbconfig.Name)
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
	email, (CASE WHEN progress=1 THEN 'NO STATUS'
		WHEN progress=2 THEN 'REJECT'
		WHEN progress=3 THEN 'APPROVED'
		WHEN progress=4 THEN 'ON PROGRESS'
		WHEN progress=5 THEN 'OFFERING - ACCEPTED'
		WHEN progress=6 THEN 'OFFERING - DECLINED'
		WHEN progress=7 THEN 'OFFERING - CANCEL'
		WHEN progress=8 THEN 'HOLD'
		WHEN progress=9 THEN 'HOLD-REJECT'
		ELSE 'CLOSED' END)as progress, phone, school, purpose, contactpersonid,
	positionapply, jobinfo, acquaintance, scheduletime,
	acquaintancename, relationship, referralName, status,
	logtimestamps, updatedDate FROM candidate WHERE formtype = 'Operational Form' ` + result

	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()

	var response []Candidate
	var nickname, phone, school, jobinfo, acquaintanceName, relationship, referralName sql.NullString
	var updated pq.NullTime
	for rows.Next() {
		var ops Candidate
		if err := rows.Scan(&ops.Fid, &ops.FfullName, &nickname,
			&ops.Femail, &ops.Fprog, &phone, &school, &ops.Fpurpose,
			&ops.Fcontactpersonid, &ops.Fposition, &jobinfo,
			&ops.Facquaintance, &ops.FscheduleTime, &acquaintanceName,
			&relationship, &referralName, &ops.Fstatus,
			&ops.Ftimestamp, &updated); err != nil {
			log.Fatal(err)
		}

		if nickname.Valid {
			temp, _ := nickname.Value()
			strTemp, _ := temp.(string)
			ops.FnickName = strTemp
		}
		if phone.Valid {
			temp, _ := phone.Value()
			strTemp, _ := temp.(string)
			ops.FphoneNumber = strTemp
		}
		if school.Valid {
			temp, _ := school.Value()
			strTemp, _ := temp.(string)
			ops.Fschool = strTemp
		}
		if jobinfo.Valid {
			temp, _ := jobinfo.Value()
			strTemp, _ := temp.(string)
			ops.Fjobinfo = strTemp
		}
		if acquaintanceName.Valid {
			temp, _ := acquaintanceName.Value()
			strTemp, _ := temp.(string)
			ops.FacquaintanceName = strTemp
		}
		if relationship.Valid {
			temp, _ := relationship.Value()
			strTemp, _ := temp.(string)
			ops.Frelationship = strTemp
		}
		if referralName.Valid {
			temp, _ := referralName.Value()
			strTemp, _ := temp.(string)
			ops.FreferralName = strTemp
		}
		if updated.Valid {
			temp, _ := updated.Value()
			strTemp, _ := temp.(time.Time)
			ops.FupdatedDate = strTemp
		}

		response = append(response, ops)

	}

	b, _ := json.MarshalIndent(response, "", "  ")
	println(string(b))
	c.Writer.Write(b)
}
