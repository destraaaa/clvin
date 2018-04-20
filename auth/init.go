package auth

import (
	"database/sql"
	"fmt"

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

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "newpassword"
	dbname   = "hris"
)

func Authentication(c *gin.Context) {
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

	c.Writer.Write([]byte("hey You are connected Auth"))
}
