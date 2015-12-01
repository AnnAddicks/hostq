package hostqueue
import (
        "fmt"
        "net/http"

        "appengine"
        "appengine/mail"
)


func sendReminder(email string, hostName string, r *http.Request) {
        c := appengine.NewContext(r)
        msg := &mail.Message{
                Sender:  "Example.com Support <support@example.com>",
                To:      []string{email},
                Subject: fmt.Sprintf("%s it is your turn to host", hostName) ,
                Body:    fmt.Sprintf(hostMessage, hostName),
        }
        if err := mail.Send(c, msg); err != nil {
                c.Errorf("Couldn't send email: %v", err)
        }
}

const hostMessage = `
%s,

It is your turn to host!  
Respond with 'yes' to host, 'no' to go to the next in line to host, or 'skip' to skip this week and host next week.
`
