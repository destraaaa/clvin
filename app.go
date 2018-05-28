package main

import (

	// "encoding/json"

	// "github.com/gin-gonic/gin"
	// "github.com/gosample/auth"
	// "github.com/gosample/chart"
	// "github.com/gosample/interviewee"

	"github.com/destraaaa/clvin/auth"
	"github.com/destraaaa/clvin/chart"
	"github.com/destraaaa/clvin/interviewee"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/interviewee/save", interviewee.WriteData)
	r.GET("/nonopsform/view", interviewee.ReadDataNon)
	r.GET("/rejected/view", interviewee.ReadRejected)
	r.GET("/opsform/view", interviewee.ReadDataOps)
	r.GET("/interviewee/email", interviewee.EmailCandidate)
	r.POST("/nonopsform/filter", chart.Filter)
	r.POST("/nonopsform/filterChart", chart.FilterChart)
	r.GET("/nonopsform/view/total", chart.Candidate)
	r.GET("/nonopsform/view/schoolpie", chart.SchoolPie)
	r.GET("/nonopsform/view/jobpie", chart.JobPie)
	r.GET("/nonopsform/view/statpie", chart.StatPie)
	r.GET("/nonopsform/view/statbar", chart.StatBar)
	r.GET("/nonopsform/view/cpbar", chart.CPBar)
	r.GET("/nonopsform/view/posbar", chart.PositionBar)
	r.GET("/authLogin/validate", auth.ShowValidate)
	r.GET("/authLogin/user", auth.ShowUser)
	r.POST("/authLogin/user", auth.WriteUser)
	r.POST("/authLogin/delete", auth.DeleteUser)

	// r.GET("/", auth.indexHandler)
	// r.GET("/login", auth.loginHandler)
	// r.GET("/auth", auth.authHandler)

	r.Run()
	// flag.Parse()
	// logging.LogInit()

	// debug := logging.Debug.Println

	// debug("app started") // message will not appear unless run with -debug switch

	// if err := agent.Listen(&agent.Options{}); err != nil {
	// 	log.Fatal(err)
	// }

	// hwm := hello.NewHelloWorldModule()
	// // interview := interviewee.WriteData()

	// http.Handle("/metrics", promhttp.Handler())

	// http.HandleFunc("/hello", hwm.SayHelloWorld)
	// go logging.StatsLog()

	// //http.HandleFunc("/interviewee/save", interviewee.WriteData)
	// // go logging.StatsLog()

	// tracer.Init(&tracer.Config{Port: 8700, Enabled: true})

	// log.Fatal(grace.Serve(":9000", nil))
}
