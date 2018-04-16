package main

import (

	// "encoding/json"

	// "github.com/tokopedia/gosample/hello"
	"github.com/gin-gonic/gin"
	"github.com/gosample/interviewee"
)

func main() {
	r := gin.Default()

	r.POST("/interviewee/save", interviewee.WriteData)
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
