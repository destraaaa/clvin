package survey

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/destraaaa/clvin/env"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	_ "github.com/lib/pq" // _ "github.com/lib/pq"
)

type ResultForm struct {
	SchemaID     int        `json:"schemaId,omitempty"`
	Questions    []string   `json:"questions,omitempty"`
	Answer       *FormJson  `json:"answer,omitempty"`
	AnswerResult []FormJson `json:"answers,omitempty"`
}

type FormJson struct {
	// Questions       []string `json:"questions"`
	json.RawMessage
}

type FormId struct {
	SchemaID []FormSchema `json:"ids,omitempty"`
	LastID   int64        `json:"lastId,omitempty"`
}
type FormSchema struct {
	ID        int       `json:"schemaId,string,omitempty"`
	Data      *FormJson `json:"schema,omitempty"`
	Questions []string  `json:"questions,omitempty"`
	Name      string    `json:"schemaName,omitempty"`
	UserID    int       `json:"userId,string,omitempty"`
}

func (j FormJson) Value() (driver.Value, error) {
	return j.MarshalJSON()
}

func (j *FormJson) Scan(src interface{}) error {
	if bytes, ok := src.([]byte); ok {
		return json.Unmarshal(bytes, j)

	}
	return errors.New(fmt.Sprint("Failed to unmarshal JSON from DB", src))
}

func WriteForm(c *gin.Context) {
	dbconfig := env.GetConfig().Database
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	var dataUser FormSchema

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

	var sqlStatement string

	if dataUser.Data == nil {
		sqlStatement = `
		UPDATE form_schema SET schema_name = $1 WHERE schema_id = $2
			`
		_, err = db.Exec(sqlStatement, dataUser.Name, dataUser.ID)

	} else {
		sqlStatement = `
			INSERT INTO form_schema(schema_id, schema, questions, is_deleted, schema_name, user_id) VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT(schema_id)
				DO UPDATE SET (schema, questions, is_deleted,schema_name, user_id) = ($2, $3, $4, $5, $6);
				`
		_, err = db.Exec(sqlStatement, dataUser.ID, dataUser.Data.RawMessage, pq.Array(dataUser.Questions), false, dataUser.Name, dataUser.UserID)
	}

	fmt.Println("sqlStatment Update2", sqlStatement)
	if err != nil {
		panic(err)
	}
}

func ReadAllForm(c *gin.Context) {
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

	sqlStatement := `SELECT schema_id, schema_name, user_id FROM form_schema WHERE is_deleted = false ORDER BY schema_id`
	sqlLastID := `SELECT max(schema_id) FROM form_schema`
	rows, err := db.Query(sqlStatement)
	rowsLastID, err := db.Query(sqlLastID)
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()

	var arrSchema []FormSchema
	var formID FormId

	for rows.Next() {
		var formSchema FormSchema
		if err := rows.Scan(&formSchema.ID, &formSchema.Name, &formSchema.UserID); err != nil {
			log.Fatal(err)
		}

		arrSchema = append(arrSchema, formSchema)
	}

	var schemaID FormId
	var lastID sql.NullInt64
	for rowsLastID.Next() {
		if err := rowsLastID.Scan(&lastID); err != nil {
			log.Fatal(err)
		}
	}

	if lastID.Valid {
		temp, _ := lastID.Value()
		intTemp, _ := temp.(int64)
		schemaID.LastID = intTemp
	}

	formID.LastID = schemaID.LastID
	formID.SchemaID = arrSchema
	fmt.Println("ini form ID", formID)
	fmt.Println("ini lastID", lastID)

	b, _ := json.MarshalIndent(formID, "", "  ")
	println(string(b))
	c.Writer.Write(b)
}

//a----------------------------------------------------------------------------

func ReadForm(c *gin.Context) {
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

	surveyID := c.Param("id")
	sqlStatement := `SELECT schema, schema_name FROM form_schema WHERE schema_id=$1 AND is_deleted=false`

	rows, err := db.Query(sqlStatement, surveyID)
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()

	// var arrSchema []FormSchema
	var formSchema FormSchema

	for rows.Next() {
		if err := rows.Scan(&formSchema.Data, &formSchema.Name); err != nil {
			log.Fatal(err)
		}

		// arrSchema = append(arrSchema, formSchema)

	}

	b, _ := json.MarshalIndent(formSchema, "", "  ")
	println(string(b))
	c.Writer.Write(b)
}

func DeleteForm(c *gin.Context) {
	dbconfig := env.GetConfig().Database
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	var formSchema FormSchema
	err := c.BindJSON(&formSchema)
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

	sqlStatement := `UPDATE form_schema SET is_deleted = true WHERE schema_id=$1 ;`
	fmt.Println(sqlStatement)
	_, err = db.Exec(sqlStatement, formSchema.ID)
	if err != nil {
		panic(err)
	}
}

func WriteResult(c *gin.Context) {
	dbconfig := env.GetConfig().Database
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

	var result ResultForm
	err := c.BindJSON(&result)
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

	sqlStatement := `INSERT INTO form_result(schema_id, answer) VALUES ((SELECT schema_id FROM form_schema WHERE schema_id = $1), $2)`
	_, err = db.Exec(sqlStatement, result.SchemaID, result.Answer.RawMessage)

	if err != nil {
		panic(err)
	}
	// for _, rs := range result {
	// 	_, err = db.Exec(sqlStatement, rs.SchemaID, rs.Label, rs.Value)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	b, _ := json.MarshalIndent("Successfully add Result ", "", "  ")
	fmt.Println(string(b))
	c.Writer.Write(b)
}

func ReadResult(c *gin.Context) {
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

	result := c.Param("id")

	fmt.Println("Successfully connected!")

	// sqlStatement := `SELECT form_result.schema_id, questions
	// FROM form_result JOIN form_schema ON form_schema.schema_id = form_result.schema_id
	// WHERE form_result.schema_id =$1`

	// sqlAnswer := `SELECT form_result.schema_id,answer, questions
	// FROM form_result JOIN form_schema ON form_schema.schema_id = form_result.schema_id
	// WHERE form_result.schema_id =$1`

	sqlStatement := `SELECT form_result.schema_id, questions
	FROM form_result JOIN form_schema ON form_schema.schema_id = form_result.schema_id
	WHERE form_result.schema_id =$1 AND is_deleted = false`

	sqlAnswer := `SELECT answer
	FROM form_result JOIN form_schema ON form_schema.schema_id = form_result.schema_id
	WHERE form_result.schema_id =$1 AND is_deleted = false`

	rows, err := db.Query(sqlStatement, result)
	rowsAnswer, err := db.Query(sqlAnswer, result)
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()
	defer rowsAnswer.Close()

	var formResult ResultForm

	for rows.Next() {
		if err := rows.Scan(&formResult.SchemaID, pq.Array(&formResult.Questions)); err != nil {
			log.Fatal(err)
		}
	}

	var resultAnswer []FormJson
	for rowsAnswer.Next() {
		var answer FormJson
		if err := rowsAnswer.Scan(&answer.RawMessage); err != nil {
			log.Fatal(err)
		}

		resultAnswer = append(resultAnswer, answer)
	}

	var joinAnswer ResultForm

	joinAnswer.SchemaID = formResult.SchemaID
	joinAnswer.AnswerResult = resultAnswer
	joinAnswer.Questions = formResult.Questions

	// var resultAnswer []ResultForm
	// for rowsAnswer.Next() {
	// 	var answer ResultForm
	// 	if err := rowsAnswer.Scan(&answer.SchemaID, &answer.Answer.RawMessage, pq.Array(&answer.Questions)); err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	resultAnswer = append(resultAnswer, answer)
	// }

	fmt.Println("resultAnswer", resultAnswer)
	fmt.Println("formResult", formResult)

	b, _ := json.MarshalIndent(joinAnswer, "", "  ")
	println(string(b))
	c.Writer.Write(b)
}
