package chart

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // _ "github.com/lib/pq"
)

type pieChart struct {
	Labels     string `json:"labels"`
	Series     int32  `json:"series"`
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

	schooldb, err := db.Query(`SELECT LOWER(school), count(school) FROM candidate GROUP BY LOWER(school)`)
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

	dashdb, err := db.Query(`SELECT 
								concat('(',round((((COUNT(jobinfo))::float/(countjob)::float)*100::float)::numeric,2)::text,'%) ', jobinfo)as percent,
								COUNT(jobinfo)
								
							FROM 
								candidate , (SELECT count(*)as countjob from candidate)as c
							GROUP BY jobinfo , countjob`)
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

	dashdb, err := db.Query(`SELECT COUNT(*) FROM candidate`)
	if err != nil {
		log.Panic(err)
	}
	defer dashdb.Close()
	// println(rows)

	var total []pieChart

	for dashdb.Next() {
		var tot pieChart
		if err := dashdb.Scan(&tot.Total); err != nil {
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

	dashdb, err := db.Query(`SELECT  count(progress) 
	from candidate where progress = 1
	UNION ALL
	SELECT  count(progress)
	from candidate where progress = 2 
	UNION ALL
	SELECT  count(progress)
	from candidate where progress = 3 
	UNION ALL
	SELECT  count(progress)
	from candidate where progress = 4`)
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

	dashdb, err := db.Query(`SELECT contactpersonid, COALESCE(sum(t.noStatus),0)as noStatus, COALESCE(sum(t.reject),0)as reject, COALESCE(sum(t.approved),0)as approved, COALESCE(sum(t.onProgress),0)as onProgress
	FROM (
	SELECT contactpersonid, (count(contactpersonid))as noStatus, (null)::bigint as reject, (null)::bigint as approved, (null)::bigint as onProgress
	from candidate where progress=1 GROUP BY contactpersonid
	UNION 
	SELECT contactpersonid, null, count(contactpersonid), null, null from candidate where progress=2 GROUP BY contactpersonid
	UNION 
	SELECT contactpersonid, null, null, count(contactpersonid), null from candidate where progress=3 GROUP BY contactpersonid
	UNION 
	SELECT contactpersonid, null, null, null, count(contactpersonid)from candidate where progress=4 GROUP BY contactpersonid
	) t  
	group by contactpersonid`)
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

	dashdb, err := db.Query(`SELECT positionapply, count(positionapply)
							 FROM candidate GROUP BY positionapply 
							 ORDER BY positionapply ASC`)
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
