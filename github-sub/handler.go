package function

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"github.com/openfaas/openfaas-cloud/sdk"
)

var db *sql.DB

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
}

func Handle(w http.ResponseWriter, r *http.Request) {
	var input []byte

	webhookSecret, _ := sdk.ReadSecret("webhook-secret")
	log.Printf("Webhook secret: %d", len(webhookSecret))

	if r.Body != nil {
		defer r.Body.Close()

		body, _ := ioutil.ReadAll(r.Body)

		input = body
	}

	err := db.Ping()
	if err != nil {
		w.WriteHeader(http.StatusOK)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ping OK"))
}
