package main

import (
	"flag"
	"log"
	"net/http"
	// "encoding/json"
	"github.com/google/gops/agent"
	// "github.com/gosample/interviewee"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	// "github.com/tokopedia/gosample/hello"
	"github.com/gosample/hello"
	"github.com/tokopedia/logging/tracer"
	"gopkg.in/tokopedia/grace.v1"
	"gopkg.in/tokopedia/logging.v1"
)

func main() {

	flag.Parse()
	logging.LogInit()

	debug := logging.Debug.Println

	debug("app started") // message will not appear unless run with -debug switch

	if err := agent.Listen(&agent.Options{}); err != nil {
		log.Fatal(err)
	}

	hwm := hello.NewHelloWorldModule()
	// interview := interviewee.RetrieveData()

	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/hello", hwm.SayHelloWorld)
	go logging.StatsLog()

	// http.HandleFunc("/interviewee/save", interview.SayHelloWorld)
	// go logging.StatsLog()

	tracer.Init(&tracer.Config{Port: 8700, Enabled: true})

	log.Fatal(grace.Serve(":9000", nil))
}
