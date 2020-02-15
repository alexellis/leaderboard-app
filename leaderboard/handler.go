package function

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"github.com/pkg/errors"

	"github.com/openfaas/openfaas-cloud/sdk"
)

var db *sql.DB
var cors string

// init establishes a persistent connection to the remote database
// the function will panic if it cannot establish a link and the
// container will restart / go into a crash/back-off loop
func init() {
	password, _ := sdk.ReadSecret("password")
	user, _ := sdk.ReadSecret("username")
	host, _ := sdk.ReadSecret("host")

	dbName := os.Getenv("postgres_db")
	port := os.Getenv("postgres_port")
	sslmode := os.Getenv("postgres_sslmode")

	connStr := "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbName + "?sslmode=" + sslmode

	var err error
	db, err = sql.Open("postgres", connStr)

	if err != nil {
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	if val, ok := os.LookupEnv("allow_cors"); ok && len(val) > 0 {
		cors = val
	}
}

// Handle a HTTP request as a middleware processor.
func Handle(w http.ResponseWriter, r *http.Request) {

	rows, getErr := db.Query(`select * from get_leaderboard();`)

	if getErr != nil {
		log.Printf("get error: %s", getErr.Error())
		http.Error(w, errors.Wrap(getErr, "unable to get from leaderboard").Error(),
			http.StatusInternalServerError)
		return
	}

	results := []Result{}
	defer rows.Close()
	for rows.Next() {
		result := Result{}
		scanErr := rows.Scan(&result.UserID, &result.UserLogin, &result.IssueComments, &result.IssuesCreated)
		if scanErr != nil {
			log.Println("scan err:", scanErr)
		}
		results = append(results, result)
	}

	if len(cors) > 0 {
		w.Header().Set("Access-Control-Allow-Origin", cors)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	res, _ := json.Marshal(results)
	w.Write(res)
}

type Result struct {
	UserID    int
	UserLogin string

	IssueComments int
	IssuesCreated int
}
