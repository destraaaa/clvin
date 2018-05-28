package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/destraaaa/clvin/env"
	"github.com/gin-gonic/gin"
)

type User struct {
	UID         int    `json:"id"`
	UName       string `json:"name"`
	UTimestamps string `json:"timestamps, omitempty"`
	UEmail      string `json:"email"`
	UDeleted    string `json:"deletes, omitempty"`
	UPic        string `json:"pic"`
	UStatus     bool   `json:"status, omitempty"`
}

func ShowValidate(c *gin.Context) {
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

	userdb, err := db.Query(`SELECT email, id FROM users`)
	if err != nil {
		log.Panic(err)
	}
	defer userdb.Close()
	// println(rows)

	var dataUser []User

	for userdb.Next() {
		var user User
		if err := userdb.Scan(&user.UEmail, &user.UID); err != nil {
			log.Fatal(err)
		}
		dataUser = append(dataUser, user)
	}

	b, _ := json.MarshalIndent(dataUser, "", "  ")
	println(string(b))
	c.Writer.Write(b)
}

func WriteUser(c *gin.Context) {
	dbconfig := env.GetConfig().Database
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	var dataUser User
	err := c.BindJSON(&dataUser)
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
	INSERT INTO users (email, name, pic, status) VALUES ($1,$2,$3, $4)`
	_, err = db.Exec(sqlStatement, dataUser.UEmail, dataUser.UName, dataUser.UPic, dataUser.UStatus)
	if err != nil {
		panic(err)
	}
}

func DeleteUser(c *gin.Context) {
	dbconfig := env.GetConfig().Database
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	var dataUser User
	err := c.BindJSON(&dataUser)
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

	sqlStatement := `UPDATE users SET status = false, deletedBy = $2 WHERE email=$1 ;`
	fmt.Println(sqlStatement)
	_, err = db.Exec(sqlStatement, dataUser.UEmail, dataUser.UDeleted)
	if err != nil {
		panic(err)
	}
}

func ShowUser(c *gin.Context) {
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

	userdb, err := db.Query(`SELECT id, name, email, pic, logtimestamps FROM users WHERE status = true`)
	if err != nil {
		log.Panic(err)
	}
	defer userdb.Close()

	var dataUser []User

	for userdb.Next() {
		var user User
		if err := userdb.Scan(&user.UID, &user.UName, &user.UEmail, &user.UPic, &user.UTimestamps); err != nil {
			log.Fatal(err)
		}
		dataUser = append(dataUser, user)
	}

	b, _ := json.MarshalIndent(dataUser, "", "  ")
	c.Writer.Write(b)
}
