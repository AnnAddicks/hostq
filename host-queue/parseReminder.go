/*
	1. Get Sender - is it one of the registered senders in queue and are they the hosting group?
*/
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
		groups, err := GetGroups(c)
		if err != nil {
			log.Fatal(err)
		}

		var g Group
		for _, group := range groups {
			for _, host := range group.Hosts {
				for _, email := range host.Emails {
					if email == from {
						g = group
						break;
					}
					continue
				}
				if g.GroupName != "" {  // cannot use (Group{}) because of []Hosts for some reason
					break;
				}
				continue

			}
			if g.GroupName != "" {
				break;
			}
			continue
		}

		if !contains(g.Hosts[0].Emails, from) {
			log.Fatal("Sent from the wrong person!\n  Sent from: %s, but expected: %s", from, strings.Join(g.Hosts[0].Emails, ", or "))
		}


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
        	// Update the order in the group
        	hosts := g.Hosts
        	currentHost := hosts[0]
        	hosts = hosts[1:]
        	hosts = append(hosts, currentHost)  //Think slices are by reference??

        	g.save(c)
        	fmt.Printf("Match Yes")
	    } else if skip.MatchString(bodyString) == true {
	    	//Respond with the current turn order for next week

	        fmt.Printf("Match Skip")
	    } else if no.MatchString(bodyString) == true {
	    	//4. If No send an email to the next in line
	    	fmt.Printf("Match No")
	    } else {
	    	c.Infof("Could not find yes/no/skip")
	    }
}

//http://stackoverflow.com/questions/10485743/contains-method-for-a-slice
func contains(slice []string, item string) bool {
    set := make(map[string]struct{}, len(slice))
    for _, s := range slice {
        set[s] = struct{}{}
    }

    _, ok := set[item] 
    return ok
}