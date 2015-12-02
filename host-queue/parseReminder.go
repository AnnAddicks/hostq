package hostqueue

import (
	"appengine"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/mail"
	"net/http"
    "regexp"
    "strings"
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
        yes, err := regexp.Compile(`yes`)
        no, err := regexp.Compile(`no`)
        skip, err := regexp.Compile(`skip`)

        body, err := ioutil.ReadAll(m.Body)
        //Convert the byte array to a lower case string
        n := bytes.IndexByte(body, 0)
        s := string(body[:n])
        bodyString := strings.ToLower(s)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s", body)


        if yes.MatchString(bodyString) == true {
        	//3. If yes or skip, respond with the current turn order, update the order in the group
        	fmt.Printf("Match Yes")
	    } else if skip.MatchString(bodyString) == true {
	    	//3. If yes or skip, respond with the current turn order
	        fmt.Printf("Match Skip")
	    } else if no.MatchString(bodyString) == true {
	    	//4. If No send an email to the next in line
	    	fmt.Printf("Match No")
	    } else {
	    	c.Infof("Could not find yes/no/skip")
	    }
}