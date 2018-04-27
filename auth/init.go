package auth

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	//"github.com/gin-gonic/contrib/sessions"
	// "golang.org/x/oauth2"
	// "golang.org/x/oauth2/google"
	// "golang.org/x/net/context"
	// "golang.org/x/oauth2"
	// "golang.org/x/oauth2/google"
	// newappengine "google.golang.org/appengine"
	// newurlfetch "google.golang.org/appengine/urlfetch"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // _ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "newpassword"
	dbname   = "hris"
)

// type User struct {
// 	Sub       string `json:"sub"`
// 	FirstName string `json:"firstName"`
// 	LastName  string `json:"LastName"`
// 	// GivenName string `json:"given_name"`
// 	// FamilyName string `json:"family_name"`
// 	// Profile string `json:"profile"`
// 	Picture       string `json:"picture"`
// 	Email         string `json:"email"`
// 	EmailVerified string `json:"email_verified"`
// 	// Gender string `json:"gender"`
// }

// type Credentials struct {
// 	Cid     string `json:"cid"`
// 	Csecret string `json:"csecret"`
// }

// var cred Credentials
// var conf *oauth2.Config
// var state string
// var store = sessions.NewCookieStore([]byte("secret"))

// func init() {
// 	file, err := ioutil.ReadFile("./creds.json")
// 	if err != nil {
// 		log.Printf("File error: %v\n", err)
// 		os.Exit(1)
// 	}
// 	json.Unmarshal(file, &cred)

// 	conf = &oauth2.Config{
// 		ClientID:     cred.Cid,
// 		ClientSecret: cred.Csecret,
// 		RedirectURL:  "http://localhost:3000",
// 		Scopes: []string{
// 			"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
// 		},
// 		Endpoint: google.Endpoint,
// 	}
// }

// func indexHandler(c *gin.Context) {
// 	c.HTML(http.StatusOK, "index.tmpl", gin.H{})
// }

// func getLoginURL(state string) string {
// 	// State can be some kind of random generated hash string.
// 	// See relevant RFC: http://tools.ietf.org/html/rfc6749#section-10.12
// 	return conf.AuthCodeURL(state)
// }

// func authHandler(c *gin.Context) {
// 	// Handle the exchange code to initiate a transport.
// 	session := sessions.Default(c)
// 	retrievedState := session.Get("state")
// 	if retrievedState != c.Query("state") {
// 		c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid session state: %s", retrievedState))
// 		return
// 	}

// 	tok, err := conf.Exchange(oauth2.NoContext, c.Query("code"))
// 	if err != nil {
// 		c.AbortWithError(http.StatusBadRequest, err)
// 		return
// 	}

// 	client := conf.Client(oauth2.NoContext, tok)
// 	email, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
// 	if err != nil {
// 		c.AbortWithError(http.StatusBadRequest, err)
// 		return
// 	}
// 	defer email.Body.Close()
// 	data, _ := ioutil.ReadAll(email.Body)
// 	log.Println("Email body: ", string(data))
// 	c.Status(http.StatusOK)
// }

// func randToken() string {
// 	b := make([]byte, 32)
// 	rand.Read(b)
// 	return base64.StdEncoding.EncodeToString(b)
// }
// func loginHandler(c *gin.Context) {
// 	state = randToken()
// 	session := sessions.Default(c)
// 	session.Set("state", state)
// 	session.Save()
// 	c.Writer.Write([]byte("<html><title>Golang Google</title> <body> <a href='" + getLoginURL() + "'><button>Login with Google!</button> </a> </body></html>"))
// }

// func Authentication(c *gin.Context) {
// 	router := gin.Default()
// 	router.Use(sessions.Sessions("goquestsession", store))
// 	router.Static("/css", "./static/css")
// 	router.Static("/img", "./static/img")
// 	router.LoadHTMLGlob("templates/*")

// 	router.GET("/", indexHandler)
// 	router.GET("/login", loginHandler)
// 	router.GET("/auth", authHandler)

// 	router.Run()
// }

type User struct {
	UId         int    `json:"id"`
	UName       string `json:"name"`
	UTimestamps string `json:"timestamps"`
	UEmail      string `json:"email"`
}

func ShowValidate(c *gin.Context) {
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

	userdb, err := db.Query(`SELECT email FROM users`)
	if err != nil {
		log.Panic(err)
	}
	defer userdb.Close()
	// println(rows)

	var dataUser []User

	for userdb.Next() {
		var user User
		if err := userdb.Scan(&user.UEmail); err != nil {
			log.Fatal(err)
		}
		dataUser = append(dataUser, user)
	}

	b, _ := json.MarshalIndent(dataUser, "", "  ")
	println(string(b))
	c.Writer.Write(b)
}

func WriteUser(c *gin.Context) {
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

	sqlStatement := `
	INSERT INTO users (email, name) VALUES ($1,$2)`
	_, err = db.Exec(sqlStatement, dataUser.UEmail, dataUser.UName)
	if err != nil {
		panic(err)
	}
}
func ShowUser(c *gin.Context) {
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

	userdb, err := db.Query(`SELECT id, name, email, logtimestamps FROM users`)
	if err != nil {
		log.Panic(err)
	}
	defer userdb.Close()
	// println(rows)

	var dataUser []User

	for userdb.Next() {
		var user User
		if err := userdb.Scan(&user.UId, &user.UName, &user.UEmail, &user.UTimestamps); err != nil {
			log.Fatal(err)
		}
		dataUser = append(dataUser, user)
	}

	b, _ := json.MarshalIndent(dataUser, "", "  ")
	println(string(b))
	c.Writer.Write(b)
}
