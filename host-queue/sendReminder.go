package hostqueue

import (
        "bytes"
        "fmt"
        "net/http"
        "strings"

        "appengine"
        "appengine/urlfetch"
        "appengine/datastore"

        "github.com/sendgrid/sendgrid-go"
)

//const emailUserName := datastore.NewQuery("emailUserName")
//const emailPassword := datastore.NewQuery("emailPassword")

/*func getEmailCreds(c appengine.Context) (emailUserName string, emailPassword string) {
        emailUserNameQuery := datastore.NewQuery("emailUserName")
        emailPasswordQuery := datastore.NewQuery("emailPassword")

}*/

func sendReminder(group Group, r *http.Request) {
        c := appengine.NewContext(r)

        creds, err := GetCreds(c)
        if err != nil { 
          panic(err)
        }

        sg := sendgrid.NewSendGridClient(creds.Username, creds.Pass)
        sg.Client = urlfetch.Client(c)
        email := group.GroupEmail
        hostName := group.Next.HostName
        c.Infof("Email: %v", email)
        c.Infof("Host Name: %v", hostName)
        

        message := sendgrid.NewMail()
        message.AddTo(email)
        subject := hostName + " it is your turn to host"
        message.SetSubject(subject)
        message.SetHTML(hostName + hostMessage )
        message.SetFrom(from)
        sg.Send(message)
        

}

func sendSkipMessage(group Group, r *http.Request) {
        c := appengine.NewContext(r)

        creds, err := GetCreds(c)
        if err != nil { 
          panic(err)
        }

        sg := sendgrid.NewSendGridClient(creds.Username, creds.Pass)
        sg.Client = urlfetch.Client(c)
        email := group.GroupEmail
        var buffer bytes.Buffer
        for _, element := range group.Hosts {
                buffer.WriteString(element.HostName)
        }
        c.Infof("buffer: %s", buffer.String())

        hosts := make([]string, len(group.Hosts))
        for i, element := range group.Hosts {
                hosts[i] = element.HostName
        }
        c.Infof("hosts: %s", strings.Join(hosts[:],","))
        c.Infof("Email: %v", email)
        
        message := sendgrid.NewMail()
        message.AddTo(email)
        message.SetSubject("See you next week")
        message.SetHTML(fmt.Sprintf(skipMessage, strings.Join(hosts[:],", ")))
        message.SetFrom(from)
        sg.Send(message)
}

/*App engine allows for environment variables, but they are stored in app.yaml
and I don't want my mail creds pushed to a public repo.  */
type EmailCreds struct {
        Username string `json:"username"`
        Pass string `json:"password"`
}

func GetCreds(ctx appengine.Context) (EmailCreds, error) {
  //q := datastore.NewQuery("EmailCreds")
  k := datastore.NewKey(ctx, "EmailCreds", "singleton_creds", 0, nil)
  var ec EmailCreds
  err := datastore.Get(ctx, k, &ec)
  if err != nil {
    return EmailCreds{}, err
  }

  return ec, nil
}

const from = "reminder@hostqueue-1146.appspotmail.com"

const hostMessage = `
,

It is your turn to host!  
Respond with 'yes' to host, 'no' to go to the next in line to host, or 'skip' for everyone to skip this week and host next week.
`

const skipMessage = `
See you next week with the following turn order:  %s
`
