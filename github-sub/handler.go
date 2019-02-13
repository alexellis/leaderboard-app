package function

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/github"
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

	if r.Body != nil {
		defer r.Body.Close()
	}

	webhookType := github.WebHookType(r)
	webhookSecret, _ := sdk.ReadSecret("webhook-secret")
	log.Printf("Webhook secret: %d", len(webhookSecret))

	body, _ := ioutil.ReadAll(r.Body)
	event, err := github.ParseWebHook(webhookType, body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	msg := ""
	if issueEvent, ok := event.(*github.IssuesEvent); ok {
		switch *issueEvent.Action {
		case "opened":
			msg = " (issue opened) by " + issueEvent.Sender.GetLogin()

			insertUser(issueEvent.Sender.GetLogin(),
				issueEvent.Sender.GetID(),
				true)
		}
	}

	if issueCommentEvent, ok := event.(*github.IssueCommentEvent); ok {
		switch *issueCommentEvent.Action {
		case "created":
			msg = " (comment created) by " + issueCommentEvent.Sender.GetLogin()

			insertUser(issueCommentEvent.Sender.GetLogin(),
				issueCommentEvent.Sender.GetID(),
				true)
		}
	}

	dbErr := db.Ping()
	if dbErr != nil {
		w.WriteHeader(http.StatusOK)
		http.Error(w, dbErr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ping OK" + msg))
}

func insertUser(login string, ID int64, track bool) error {
	_, err := db.Query(`insert into users`+
		` (user_id, user_login, track) values ($1, $2, $3);`,
		login, int(ID), track)

	return err
}
