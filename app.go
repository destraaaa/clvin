package main

import (
	"github.com/destraaaa/clvin/auth"
	"github.com/destraaaa/clvin/chart"
	"github.com/destraaaa/clvin/interviewee"
	"github.com/destraaaa/clvin/survey"
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
	r.POST("/form/update", interviewee.UpdateCandidate)
	r.GET("/nonopsform/view/total", chart.Candidate)
	r.GET("/nonopsform/view/schoolpie", chart.SchoolPie)
	r.GET("/nonopsform/view/jobpie", chart.JobPie)
	r.GET("/nonopsform/view/statpie", chart.StatPie)
	r.GET("/nonopsform/view/statbar", chart.StatBar)
	r.GET("/nonopsform/view/cpbar", chart.CPBar)
	r.GET("/nonopsform/view/posbar", chart.PositionBar)
	r.GET("/form/schoolRegist", chart.SchoolRegist)
	r.GET("/authLogin/validate", auth.ShowValidate)
	r.GET("/authLogin/user", auth.ShowUser)
	r.POST("/authLogin/user", auth.WriteUser)
	r.POST("/authLogin/delete", auth.DeleteUser)
	r.POST("/survey/form/save", survey.WriteForm)
	r.GET("/survey/form/view", survey.ReadAllForm)
	r.GET("/survey/form/view/:id", survey.ReadForm)
	r.POST("/survey/form/delete", survey.DeleteForm)
	r.POST("/survey/result/save", survey.WriteResult)
	r.GET("/survey/result/view/:id", survey.ReadResult)

	r.Run()
}
