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
			login := issueEvent.Sender.GetLogin()
			id := issueEvent.Sender.GetID()
			msg = " (issue opened) by " + login
			owner := *issueEvent.Repo.GetOwner().Login
			repo := issueEvent.Repo.GetName()
			activityType := "issue_created"
			insertErr := insertUser(login, id, true)
			if insertErr != nil {
				log.Printf("%s\n", insertErr.Error())
			}

			//insert into activity (id,user_id,activity_type,activity_date,owner,repo) values (DEFAULT,653013,'issue_created','2019-02-13 07:44:00','openfaas','org-tester');
			activityErr := insertActivity(id, activityType, owner, repo)
			if activityErr != nil {
				log.Printf("%s\n", activityErr.Error())
			}
		}
	}

	if issueCommentEvent, ok := event.(*github.IssueCommentEvent); ok {
		switch *issueCommentEvent.Action {
		case "created":

			msg = " (comment created) by " + issueCommentEvent.Sender.GetLogin()
			login := issueCommentEvent.Sender.GetLogin()
			id := issueCommentEvent.Sender.GetID()
			owner := *issueCommentEvent.Repo.GetOwner().Login
			repo := issueCommentEvent.Repo.GetName()
			activityType := "issue_comment"

			insertErr := insertUser(login, id, true)

			if insertErr != nil {
				log.Printf("%s\n", insertErr.Error())
			}

			activityErr := insertActivity(id, activityType, owner, repo)
			if activityErr != nil {
				log.Printf("%s\n", activityErr.Error())
			}
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
	_, err := db.Query(`insert into users (user_id, user_login, track, created_at) values ($1, $2, $3, now());`,
		ID, login, track)

	return err
}

//insert into activity (id,user_id,activity_type,activity_date,owner,repo) values (DEFAULT,653013,'issue_created','2019-02-13 07:44:00','openfaas','org-tester');
func insertActivity(loginID int64, activityType, owner, repo string) error {
	_, err := db.Query(`insert into activity (id,user_id,activity_type,activity_date,owner,repo) values (DEFAULT,$1, $2, now(), $4, $5);`,
		loginID, activityType, owner, repo)

	return err
}
