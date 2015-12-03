package hostqueue

import (
        "bytes"
        "fmt"
        "net/http"
        "strings"

        "appengine"
        "appengine/mail"
)


func sendReminder(group Group, r *http.Request) {
        c := appengine.NewContext(r)
        email := group.GroupEmail
        hostName := group.Hosts[0].HostName

        msg := &mail.Message{
                Sender:  "reminder@hostqueue-1146.appspotmail.com Support <reminder@hostqueue-1146.appspotmail.com>",
                To:      []string{email},
                Subject: fmt.Sprintf("%s it is your turn to host", hostName) ,
                Body:    fmt.Sprintf(hostMessage, hostName),
        }
        if err := mail.Send(c, msg); err != nil {
                c.Errorf("Couldn't send email: %v", err)
        }
}

func sendSkipMessage(group Group, r *http.Request) {
        c := appengine.NewContext(r)
        email := group.GroupEmail
        
        var buffer bytes.Buffer
        for _, element := range group.Hosts {
                buffer.WriteString(element.HostName)
        }
        c.Infof("buffer: %s", buffer.String())

        hosts := []string {}
        for i, element := range group.Hosts {
                hosts[i] = element.HostName
        }
        c.Infof("hosts: %s", strings.Join(hosts[:],","))


        msg := &mail.Message{
                Sender:  "reminder@hostqueue-1146.appspotmail.com Support <reminder@hostqueue-1146.appspotmail.com>",
                To:      []string{email},
                Subject: fmt.Sprintf("See you next week") ,
                Body:    fmt.Sprintf(skipMessage, strings.Join(hosts[:],",")),
        }
        if err := mail.Send(c, msg); err != nil {
                c.Errorf("Couldn't send email: %v", err)
        }
}

const hostMessage = `
%s,

It is your turn to host!  
Respond with 'yes' to host, 'no' to go to the next in line to host, or 'skip' for everyone to skip this week and host next week.
`

const skipMessage = `
See you next week with the following turn order:  %s
`
