package hostqueue

import (
	"appengine"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/mail"
	"net/http"
)

func init() {
        http.HandleFunc("/_ah/mail/", incomingMail)
}

func incomingMail(w http.ResponseWriter, r *http.Request) {
        //Sample from https://cloud.google.com/appengine/docs/go/mail/
       
        c := appengine.NewContext(r)
        defer r.Body.Close()
        var b bytes.Buffer
        if _, err := b.ReadFrom(r.Body); err != nil {
                c.Errorf("Error reading body: %v", err)
                return
        }
        c.Infof("Received mail: %v", b)


        //Steps
        //1. Get Sender - is it one of the registered senders in queue and are they the hosting group?
        m, err := mail.ReadMessage(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		header := m.Header
		from := header.Get("From")
		c.Infof("Email replied from: %s", from)
		//***** Check if they should be responding or if someone is being snarky. ************

        //2. Look for Yes/No/Skip
        body, err := ioutil.ReadAll(m.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s", body)


        //3. If yes or skip, respond with the current turn order, update the order in the group
        //4. If No send an email to the next in line
}