package chart

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // _ "github.com/lib/pq"
)

type filter struct {
	Years   string `json:"year"`
	Type    string `json:"type"`
	Quarter string `json:"quarter"`
	Month   string `json:"month"`
	Daily   string `json:"daily"`
}

type pieChart struct {
	Labels     string `json:"labels,omitempty"`
	Series     int32  `json:"series,omitempty"`
	NoStatus   int32  `json:"nostatus"`
	Reject     int32  `json:"reject"`
	Approved   int32  `json:"approved"`
	OnProgress int32  `json:"onprogress"`
	Total      int32  `json:"total"`
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

var Fil filter

func Filter(c *gin.Context) {
	// fmt.Printf([]byte(r.Body))

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	err := c.BindJSON(&Fil)
	if err != nil {
		fmt.Println("Error Binding JSON")
		fmt.Println(err.Error())
	}

	// fmt.Println(fil.Daily, fil.Month, fil.Quarter, fil.Type, fil.Years)

}

//FilteringValidate
func Filtering() string {
	var where []string
	var temp string

	if Fil.Years != "all" && Fil.Years != "" {
		temp = "AND EXTRACT(YEAR FROM logtimestamps)= " + Fil.Years
		where = append(where, temp)
	}
	if Fil.Type != "all" && Fil.Type != "" {
		temp = "AND formtype = '" + Fil.Type + "'"
		where = append(where, temp)
	}
	if Fil.Quarter != "all" && Fil.Quarter != "" {
		temp = "AND EXTRACT(QUARTER FROM logtimestamps)= " + Fil.Quarter
		where = append(where, temp)
	}
	if Fil.Month != "all" && Fil.Month != "" {
		temp = "AND EXTRACT(MONTH FROM logtimestamps)= " + Fil.Month
		where = append(where, temp)
	}
	if Fil.Daily == "week" {
		temp = "AND logtimestamps >= now() - interval '1 week' "
		where = append(where, temp)
	}
	if Fil.Daily == "day" {
		temp = "AND EXTRACT(DAY FROM logtimestamps)= EXTRACT(DAY FROM CURRENT_TIMESTAMP)"
		where = append(where, temp)
	}

	result := strings.Join(where, " ")

	return result

}

func SchoolPie(c *gin.Context) {
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

	// var fil filter
	// errs := c.BindJSON(&fil)
	// if errs != nil {
	// 	fmt.Println("Error Binding JSON")
	// 	fmt.Println(err.Error())
	// }

	result := Filtering()

	// schooldb, err := db.Query("SELECT LOWER(school), count(school) FROM candidate WHERE school=school AND EXTRACT(YEAR FROM logtimestamps)= 2018 AND formtype = 'Non Operational Form' AND EXTRACT(QUARTER FROM logtimestamps)= 2 AND EXTRACT(MONTH FROM logtimestamps)= 5 GROUP BY LOWER(school)")

	sqlstatment := "SELECT LOWER(school), count(school) FROM candidate WHERE school=school " + result + " GROUP BY LOWER(school)"
	// fmt.Println(sqlstatment)
	schooldb, err := db.Query(sqlstatment)
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

func JobPie(c *gin.Context) {
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

	result := Filtering()

	sqlStatment := `SELECT 
					concat('(',round((((COUNT(jobinfo))::float/(countjob)::float)*100::float)::numeric,2)::text,'%) ', jobinfo)as percent,
					COUNT(jobinfo)
					
					FROM 
						candidate , (SELECT count(*)as countjob from candidate)as c
					WHERE
						jobinfo=jobinfo ` + result + ` GROUP BY jobinfo , countjob`

	dashdb, err := db.Query(sqlStatment)
	if err != nil {
		log.Panic(err)
	}
	defer dashdb.Close()
	// println(rows)

	var dataJob []pieChart

	for dashdb.Next() {
		var job pieChart
		if err := dashdb.Scan(&job.Labels, &job.Series); err != nil {
			log.Fatal(err)
		}
		dataJob = append(dataJob, job)
	}

	b, _ := json.MarshalIndent(dataJob, "", "  ")
	println(string(b))
	c.Writer.Write(b)
}

func StatPie(c *gin.Context) {
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

	result := Filtering()

	sqlStatment := `SELECT concat('(',round((((COUNT(progress))::float/(countProgress)::float)*100::float)::numeric,2)::text,'%) ', 
					(CASE WHEN progress=1 THEN 'no status'
						WHEN progress=2 THEN 'reject'
						WHEN progress=3 THEN 'approved'
						ELSE 'on progress' END))as percent, COUNT(progress)
					FROM candidate , (SELECT count(*)as countProgress from candidate)as c
					WHERE progress=progress ` + result + ` GROUP BY progress , countProgress`

	dashdb, err := db.Query(sqlStatment)
	if err != nil {
		log.Panic(err)
	}
	defer dashdb.Close()
	// println(rows)

	var dataStat []pieChart

	for dashdb.Next() {
		var stat pieChart
		if err := dashdb.Scan(&stat.Labels, &stat.Series); err != nil {
			log.Fatal(err)
		}
		dataStat = append(dataStat, stat)
	}

	b, _ := json.MarshalIndent(dataStat, "", "  ")
	println(string(b))
	c.Writer.Write(b)
}

func Candidate(c *gin.Context) {
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

	result := Filtering()
	sqlStatment := `SELECT count(*), (SELECT count(*) FROM candidate WHERE progress=2 ` + result + ` ),
	(SELECT count(*) FROM candidate WHERE progress=3 ` + result + ` ), (SELECT count(*)
	FROM candidate WHERE progress=4 ` + result + ` ) FROM candidate WHERE progress = progress ` + result

	dashdb, err := db.Query(sqlStatment)
	if err != nil {
		log.Panic(err)
	}
	defer dashdb.Close()
	// println(rows)

	var total []pieChart

	for dashdb.Next() {
		var tot pieChart
		if err := dashdb.Scan(&tot.Total, &tot.Reject, &tot.Approved, &tot.OnProgress); err != nil {
			log.Fatal(err)
		}
		total = append(total, tot)
	}

	b, _ := json.MarshalIndent(total, "", "  ")
	println(string(b))
	c.Writer.Write(b)
}

func StatBar(c *gin.Context) {
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

	result := Filtering()

	sqlStatment := `SELECT  count(progress) 
					from candidate where progress = 1 ` + result +
		` UNION ALL
					SELECT  count(progress)
					from candidate where progress = 2 ` + result +
		` UNION ALL
					SELECT  count(progress)
					from candidate where progress = 3 ` + result +
		` UNION ALL
					SELECT  count(progress)
					from candidate where progress = 4 ` + result

	dashdb, err := db.Query(sqlStatment)
	if err != nil {
		log.Panic(err)
	}
	defer dashdb.Close()
	// println(rows)

	var dataProgress []pieChart

	for dashdb.Next() {
		var progress pieChart
		if err := dashdb.Scan(&progress.Series); err != nil {
			log.Fatal(err)
		}
		dataProgress = append(dataProgress, progress)
	}

	b, _ := json.MarshalIndent(dataProgress, "", "  ")
	println(string(b))
	c.Writer.Write(b)
}

func CPBar(c *gin.Context) {
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

	result := Filtering()

	sqlStatment := `SELECT t.contactpersonid, COALESCE(sum(t.noStatus),0)as noStatus, COALESCE(sum(t.reject),0)as reject, COALESCE(sum(t.approved),0)as approved, COALESCE(sum(t.onProgress),0)as onProgress
					FROM (
					SELECT contactpersonid, (count(contactpersonid))as noStatus, (null)::bigint as reject, (null)::bigint as approved, (null)::bigint as onProgress
					from candidate where progress=1 ` + result + ` GROUP BY contactpersonid
					UNION 
					SELECT contactpersonid, null, count(contactpersonid), null, null from candidate where progress=2 ` + result + ` GROUP BY contactpersonid
					UNION 
					SELECT contactpersonid, null, null, count(contactpersonid), null from candidate where progress=3 ` + result + ` GROUP BY contactpersonid
					UNION 
					SELECT contactpersonid, null, null, null, count(contactpersonid)from candidate where progress=4 ` + result + ` GROUP BY contactpersonid
					) t GROUP BY t.contactpersonid`

	dashdb, err := db.Query(sqlStatment)
	if err != nil {
		log.Panic(err)
	}
	defer dashdb.Close()
	// println(rows)

	var Cbar []pieChart

	for dashdb.Next() {
		var bar pieChart
		if err := dashdb.Scan(&bar.Labels, &bar.NoStatus, &bar.Reject, &bar.Approved,
			&bar.OnProgress); err != nil {
			log.Fatal(err)
		}
		Cbar = append(Cbar, bar)
	}

	b, _ := json.MarshalIndent(Cbar, "", "  ")
	println(string(b))
	c.Writer.Write(b)
}

func PositionBar(c *gin.Context) {
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

	result := Filtering()

	sqlStatment := `SELECT positionapply, count(positionapply)
					FROM candidate WHERE positionapply = positionapply ` + result + ` GROUP BY positionapply 
					ORDER BY positionapply ASC`
	dashdb, err := db.Query(sqlStatment)

	if err != nil {
		log.Panic(err)
	}
	defer dashdb.Close()
	// println(rows)

	var dataPos []pieChart

	for dashdb.Next() {
		var pos pieChart
		if err := dashdb.Scan(&pos.Labels, &pos.Series); err != nil {
			log.Fatal(err)
		}
		dataPos = append(dataPos, pos)
	}

	b, _ := json.MarshalIndent(dataPos, "", "  ")
	println(string(b))
	c.Writer.Write(b)
}
