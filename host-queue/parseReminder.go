/*	
	Script to handle an email response to the reminder.  
	Steps:
		1. Get Sender - is it one of the registered senders in queue and are they the hosting group?
		2. Check if they should be responding or if someone is being snarky.
		3. Look for Yes/No/Skip
			3a. Yes - update the order in the group
			3b. No - send an email to the next in line, update group
			3c. Skip - respond with the current turn order for next week
	Notes: Super procedural right now.  I need to clean up the code once I have it working! 

*/
package hostqueue

import (
	"appengine"
	"bytes"
	"github.com/go-martini/martini"
	"fmt"
	"io/ioutil"
	"log"
	"net/mail"
	"net/http"
    "regexp"
    "strings"
)

func init() {
        
}

func IncomingMail(c martini.Context, w http.ResponseWriter, r *http.Request) {
        //Sample from https://cloud.google.com/appengine/docs/go/mail/
       
        ctx := appengine.NewContext(r)
        defer r.Body.Close()
        var b bytes.Buffer
        if _, err := b.ReadFrom(r.Body); err != nil {
                ctx.Errorf("Error reading body: %v", err)
                return
        }
        ctx.Infof("Received mail: %v", b)

        //Get Sender - is it one of the registered senders in queue and are they the hosting group?
        m, err := mail.ReadMessage(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		header := m.Header
		from := header.Get("From")
		ctx.Infof("Email replied from: %s", from)
		
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

		if !contains(g.Next.Emails, from) {
			log.Fatal("Sent from the wrong person!\n  Sent from: %s, but expected: %s", from, strings.Join(g.Hosts[0].Emails, ", or "))
		}


        //Look for Yes/No/Skip
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


		//TODO:  Buggy Logic, test when working!
        if yes.MatchString(bodyString) == true {
        	// Update the order in the group
        	hosts := g.Hosts
        	currentHost := hosts[0]
        	hosts = hosts[1:]
        	hosts = append(hosts, currentHost)  //Think slices are by reference??
        	g.Next = hosts[0]
        	g.save(c)
        	fmt.Printf("Match Yes")
	    } else if skip.MatchString(bodyString) == true {
	    	//Respond with the current turn order for next week
	    	sendSkipMessage(g, r)
	        fmt.Printf("Match Skip")
	    } else if no.MatchString(bodyString) == true {
	    	//Send an email to the next in line
	    	hosts := g.Hosts
	    	currentIndex := SliceIndex(len(hosts), func(i int) bool { return contains(hosts[i].Emails, from) }) 
	    	if(currentIndex < (len(hosts) - 1)) {
	    		g.Next = hosts[currentIndex + 1]
	    	} else {
	    		g.Next = hosts[0]
	    	}

	    	g.save(c)
	    	sendReminder(g, r)
	    	fmt.Printf("Match No")
	    } else {
	    	ctx.Infof("Could not find yes/no/skip")
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

//http://stackoverflow.com/questions/8307478/go-how-to-find-out-element-position-in-slice
func SliceIndex(limit int, predicate func(i int) bool) int {
    for i := 0; i < limit; i++ {
        if predicate(i) {
            return i
        }
    }
    return -1
}