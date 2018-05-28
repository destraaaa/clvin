package chart

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/destraaaa/clvin/env"
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

type filterChart struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type pieChart struct {
	Labels           string `json:"labels,omitempty"`
	Series           int32  `json:"series,omitempty"`
	Total            int32  `json:"total"`
	NoStatus         int32  `json:"nostatus"`
	Reject           int32  `json:"reject"`
	Approved         int32  `json:"approved"`
	OnProgress       int32  `json:"onprogress"`
	OfferingAccepted int32  `json:"offeringAccepted"`
	OfferingCancel   int32  `json:"offeringCancel"`
	OfferingDeclined int32  `json:"offeringDeclined"`
	Holds            int32  `json:"holds"`
	HoldsReject      int32  `json:"holdsReject"`
	Closed           int32  `json:"closed"`
}

var Fil filter
var FilChart filterChart

func Filter(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	err := c.BindJSON(&Fil)
	if err != nil {
		fmt.Println("Error Binding JSON")
		fmt.Println(err.Error())
	}
}

func FilterChart(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	err := c.BindJSON(&FilChart)
	if err != nil {
		fmt.Println("Error Binding JSON")
		fmt.Println(err.Error())
	}

}

func FilteringChart() string {
	var temp string
	if FilChart.Type != "all" {
		temp = " LIMIT " + FilChart.Value
	} else {
		temp = ""
	}

	return temp
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
		temp = "AND EXTRACT(DAY FROM logtimestamps)= EXTRACT(DAY FROM CURRENT_TIMESTAMP) AND EXTRACT(MONTH FROM logtimestamps)= EXTRACT(MONTH FROM CURRENT_TIMESTAMP) AND EXTRACT(YEAR FROM logtimestamps)= EXTRACT(YEAR FROM CURRENT_TIMESTAMP)"
		where = append(where, temp)
	}

	result := strings.Join(where, " ")

	return result

}

func SchoolPie(c *gin.Context) {
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

	result := Filtering()
	limit := FilteringChart()
	sqlstatment := ""

	if FilChart.Type == "school" {
		sqlstatment = "SELECT LOWER(school), count(school) FROM candidate WHERE school=school " + result + " GROUP BY LOWER(school)" + limit
	} else {
		sqlstatment = "SELECT LOWER(school), count(school) FROM candidate WHERE school=school " + result + " GROUP BY LOWER(school)"
	}

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

	result := Filtering()

	sqlStatment := `SELECT 
					concat('(',round((((COUNT(jobinfo))::float/(countjob)::float)*100::float)::numeric,2)::text,'%) ', jobinfo)as percent,
					COUNT(jobinfo)
					
					FROM 
						candidate , (SELECT count(*)as countjob from candidate)as c
					WHERE
						jobinfo=jobinfo ` + result + ` GROUP BY jobinfo , countjob`

	println(sqlStatment)
	dashdb, err := db.Query(sqlStatment)
	if err != nil {
		log.Panic(err)
	}
	defer dashdb.Close()

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

	result := Filtering()

	sqlStatment := `SELECT concat('(',round((((COUNT(progress))::float/(countProgress)::float)*100::float)::numeric,2)::text,'%) ', 
					(CASE WHEN progress=1 THEN 'NO STATUS'
						WHEN progress=2 THEN 'REJECT'
						WHEN progress=3 THEN 'APPROVED'
						WHEN progress=4 THEN 'ON PROGRESS'
						WHEN progress=5 THEN 'OFFERING - ACCEPTED'
						WHEN progress=6 THEN 'OFFERING - DECLINED'
						WHEN progress=7 THEN 'OFFERING - CANCEL'
						WHEN progress=8 THEN 'HOLD'
						WHEN progress=9 THEN 'HOLD-REJECT'
						ELSE 'CLOSED' END))as percent, COUNT(progress)
					FROM candidate , (SELECT count(*)as countProgress from candidate)as c
					WHERE progress=progress ` + result + ` GROUP BY progress , countProgress`

	dashdb, err := db.Query(sqlStatment)
	if err != nil {
		log.Panic(err)
	}
	defer dashdb.Close()

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

	result := Filtering()
	sqlStatment := `SELECT count(*), (SELECT count(*) FROM candidate WHERE progress=2 ` + result + ` ),
	(SELECT count(*) FROM candidate WHERE progress=3 ` + result + ` ), (SELECT count(*)
	FROM candidate WHERE progress=4 ` + result + ` )
	 FROM candidate WHERE progress = progress ` + result

	dashdb, err := db.Query(sqlStatment)
	if err != nil {
		log.Panic(err)
	}
	defer dashdb.Close()

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
					from candidate where progress = 4 ` + result +
		` UNION ALL
					SELECT  count(progress)
					from candidate where progress = 5 ` + result +
		` UNION ALL
					SELECT  count(progress)
					from candidate where progress = 6 ` + result +
		` UNION ALL
					SELECT  count(progress)
					from candidate where progress = 7 ` + result +
		` UNION ALL
					SELECT  count(progress)
					from candidate where progress = 8 ` + result +
		` UNION ALL
					SELECT  count(progress)
					from candidate where progress = 9 ` + result +
		` UNION ALL
					SELECT  count(progress)
					from candidate where progress = 10 ` + result

	dashdb, err := db.Query(sqlStatment)
	if err != nil {
		log.Panic(err)
	}
	defer dashdb.Close()
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

	result := Filtering()

	sqlStatment := `SELECT t.contactpersonid, COALESCE(sum(t.noStatus),0)as noStatus, COALESCE(sum(t.reject),0)as reject, COALESCE(sum(t.approved),0)as approved, COALESCE(sum(t.onProgress),0)as onProgress, COALESCE(sum(t.OfferingAccepted),0)as OfferingAccepted, COALESCE(sum(t.OfferingDeclined),0)as OfferingDeclined, COALESCE(sum(t.OfferingCancel),0)as OfferingCancel, COALESCE(sum(t.holds),0)as holds, COALESCE(sum(t.holdsReject),0)as holdsReject, COALESCE(sum(t.closed),0)as closed
					FROM (
					SELECT contactpersonid, (count(contactpersonid))as noStatus, (null)::bigint as reject, (null)::bigint as approved, (null)::bigint as onProgress, (null)::bigint as OfferingAccepted, (null)::bigint as OfferingDeclined, (null)::bigint as OfferingCancel, (null)::bigint as holds, (null)::bigint as holdsReject, (null)::bigint as closed
					from candidate where progress=1 ` + result + ` GROUP BY contactpersonid
					UNION 
					SELECT contactpersonid, null, count(contactpersonid), null, null, null, null, null, null, null, null from candidate where progress=2 ` + result + ` GROUP BY contactpersonid
					UNION 
					SELECT contactpersonid, null, null, count(contactpersonid), null, null, null, null, null, null, null from candidate where progress=3 ` + result + ` GROUP BY contactpersonid
					UNION 
					SELECT contactpersonid, null, null, null, count(contactpersonid), null, null, null, null, null, null from candidate where progress=4 ` + result + ` GROUP BY contactpersonid
					UNION 
					SELECT contactpersonid, null, null, null, null, count(contactpersonid), null, null, null, null, null from candidate where progress=5 ` + result + ` GROUP BY contactpersonid
					UNION 
					SELECT contactpersonid, null, null, null, null, null, count(contactpersonid), null, null, null, null from candidate where progress=6 ` + result + ` GROUP BY contactpersonid
					UNION 
					SELECT contactpersonid, null, null, null, null, null, null, count(contactpersonid), null, null, null from candidate where progress=7 ` + result + ` GROUP BY contactpersonid
					UNION 
					SELECT contactpersonid, null, null, null, null, null, null, null, count(contactpersonid), null, null from candidate where progress=8 ` + result + ` GROUP BY contactpersonid
					UNION 
					SELECT contactpersonid, null, null, null, null, null, null, null, null, count(contactpersonid), null from candidate where progress=9 ` + result + ` GROUP BY contactpersonid
					UNION 
					SELECT contactpersonid, null, null, null, null, null, null, null, null, null, count(contactpersonid) from candidate where progress=10 ` + result + ` GROUP BY contactpersonid
					) t GROUP BY t.contactpersonid`

	fmt.Println(sqlStatment)
	dashdb, err := db.Query(sqlStatment)
	if err != nil {
		log.Panic(err)
	}
	defer dashdb.Close()

	var Cbar []pieChart

	for dashdb.Next() {
		var bar pieChart
		if err := dashdb.Scan(&bar.Labels, &bar.NoStatus, &bar.Reject, &bar.Approved,
			&bar.OnProgress, &bar.OfferingAccepted, &bar.OfferingDeclined, &bar.OfferingCancel,
			&bar.Holds, &bar.HoldsReject, &bar.Closed); err != nil {
			log.Fatal(err)
		}
		Cbar = append(Cbar, bar)
	}

	b, _ := json.MarshalIndent(Cbar, "", "  ")
	println(string(b))
	c.Writer.Write(b)
}

func PositionBar(c *gin.Context) {
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

	result := Filtering()
	limit := FilteringChart()
	sqlStatment := ""

	if FilChart.Type == "position" {
		sqlStatment = `SELECT positionapply, count(positionapply)
		FROM candidate WHERE positionapply = positionapply ` + result + ` GROUP BY positionapply 
		ORDER BY positionapply ASC` + limit
	} else {
		sqlStatment = `SELECT positionapply, count(positionapply)
					FROM candidate WHERE positionapply = positionapply ` + result + ` GROUP BY positionapply 
					ORDER BY positionapply ASC`
	}

	dashdb, err := db.Query(sqlStatment)

	if err != nil {
		log.Panic(err)
	}
	defer dashdb.Close()

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
