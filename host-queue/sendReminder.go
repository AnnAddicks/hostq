package hostqueue

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"appengine"
	"appengine/datastore"
	"appengine/urlfetch"

	"github.com/sendgrid/sendgrid-go"
)

func sendReminder(group Group, r *http.Request) {
	c := appengine.NewContext(r)

	creds, err := GetCreds(c)
	if err != nil {
		panic(err)
	}

	sg := sendgrid.NewSendGridClient(creds.Username, creds.Pass)
	sg.Client = urlfetch.Client(c)
	email := group.GroupEmail
	hostName := group.Hosts[group.Next].HostName
	c.Infof("Email: %v", email)
	c.Infof("Host Name: %v", hostName)

	subject := "This weeks hosting reminder for " + group.GroupName
	html := hostName + hostMessage
	sendEmail(email, subject, html, r)
}

func sendSkipMessage(group Group, r *http.Request) {
	c := appengine.NewContext(r)

	email := group.GroupEmail
	var buffer bytes.Buffer
	for _, element := range group.Hosts {
		buffer.WriteString(element.HostName)
	}

	hosts := make([]string, len(group.Hosts))
	for i, element := range group.Hosts {
		hosts[i] = element.HostName
	}
	c.Infof("hosts: %s", strings.Join(hosts[:], ","))
	c.Infof("Email: %v", email)

	subject := "Re: This weeks hosting reminder for " + group.GroupName
	html := fmt.Sprintf(skipMessage, strings.Join(hosts[:], ", "))
	sendEmail(email, subject, html, r)
}

func sendHostConfirmedMessage(group Group, r *http.Request) {
	c := appengine.NewContext(r)

	email := group.GroupEmail
	var buffer bytes.Buffer
	for _, element := range group.Hosts {
		buffer.WriteString(element.HostName)
	}

	hosts := make([]string, len(group.Hosts))
	for i, element := range group.Hosts {
		hosts[i] = element.HostName
	}
	c.Infof("hosts: %s", strings.Join(hosts[:], ","))
	c.Infof("Email: %v", email)

	html := fmt.Sprintf(confirmedMessage, group.Hosts[len(group.Hosts)-1].HostName, strings.Join(hosts[:], ", "))
	sendEmail(email, "Re: This weeks hosting reminder for "+group.GroupName, html, r)

}

func sendEmail(email string, subject string, html string, r *http.Request) {
	c := appengine.NewContext(r)

	creds, err := GetCreds(c)
	if err != nil {
		CreateGithubIssueAndPanic(err, r)
	}

	sg := sendgrid.NewSendGridClient(creds.Username, creds.Pass)
	sg.Client = urlfetch.Client(c)

	message := sendgrid.NewMail()
	message.AddTo(email)
	message.SetSubject(subject)
	message.SetHTML(html)
	message.SetFrom(from)

	err = sg.Send(message)
	if err != nil {
		//TODO:  autocreate an issue in github if this is unique in the last month
		c.Infof("Message: %v", message)
		CreateGithubIssueAndPanic(err, r)
	}
}

/*App engine allows for environment variables, but they are stored in app.yaml
and I don't want my mail creds pushed to a public repo.  */
type EmailCreds struct {
	Username string `json:"username"`
	Pass     string `json:"password"`
}

func GetCreds(ctx appengine.Context) (EmailCreds, error) {
	k := datastore.NewKey(ctx, "EmailCreds", "singleton_creds", 0, nil)
	var ec EmailCreds
	err := datastore.Get(ctx, k, &ec)
	if err != nil {
		return EmailCreds{}, err
	}

	return ec, nil
}

const from = "reminder@hostqueue-1146.appspotmail.com"

const hostMessage = ` it is your turn to host!  
Respond with 'yes' to host, 'no' to go to the next in line to host, or 'skip' for everyone to skip this week and host next week.

To get the live rotation status go here:  http://hostqueue-1146.appspot.com/group/status/a2b57006-7df8-444a-842a-dd366cd228bd
`

const skipMessage = `See you next week with the following turn order:  %s
To get the live rotation status go here:  http://hostqueue-1146.appspot.com/group/status/a2b57006-7df8-444a-842a-dd366cd228bd
`

const confirmedMessage = `The %s has agreed to host this week.  The rotation for next week will be: %s  

To get the live rotation status go here:  http://hostqueue-1146.appspot.com/group/status/a2b57006-7df8-444a-842a-dd366cd228bd
`
