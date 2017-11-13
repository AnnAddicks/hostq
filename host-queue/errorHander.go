/* Auto Create Github issues on crashes.
 Sample curl to create an issue:
curl --user "caspyin" -X POST --data
  '{"description":"Created via API","public":"true","files":{"file1.txt":{"content":"Demo"}}'
  https://api.github.com/gists
*/
package hostqueue

import (
	"appengine"
	"appengine/datastore"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const githubURI = "https://api.github.com/repos/AnnAddicks/hostq/issues"

type Issue struct {
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Labels []string `json:"labels"`
}

type GithubCreds struct {
	Username string `json:"username"`
	Password string `json:"username"`
}

/*
	Storing creds in the datastore because app engine puts environment
	variables in app.yaml.  TODO: Cache these
*/
func GetGithubCreds(ctx appengine.Context) (GithubCreds, error) {
	q := datastore.NewQuery("GithubCreds").Limit(1)

	var creds []GithubCreds
	var cred GithubCreds
	if _, err := q.GetAll(ctx, &creds); err != nil {
		return cred, err
	}

	return creds[0], nil
}

func CreateGithubIssueAndPanic(originalErr error, r *http.Request) {
	/* TODO check the datastore first if the error has happened in the last
	30 or so days */

	var issue Issue
	issue.Title = "Auto Created issue"
	issue.Body = fmt.Sprintf("%v", originalErr)
	issue.Labels = []string{"AutoGenerated Error"}

	jsonStr, _ := json.Marshal(issue)

	req, _ := http.NewRequest("POST", githubURI, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	ctx := appengine.NewContext(r)
	githubCreds, _ := GetGithubCreds(ctx)
	req.SetBasicAuth(githubCreds.Username, githubCreds.Password)
	handleURLFetch(req, r)

	panic(originalErr)
}