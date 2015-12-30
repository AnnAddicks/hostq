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
	and wow! look at the use of 3 different logging!

*/
package hostqueue

import (
	"appengine"
	"io/ioutil"
	"log"
	"net/mail"
	"net/http"
    "regexp"
    "strings"
)

func init() {
        
}

func IncomingMail(w http.ResponseWriter, r *http.Request) {
        ctx := appengine.NewContext(r)
        defer r.Body.Close()

        //Get Sender - is it one of the registered senders in queue and are they the hosting group?
        m, err := mail.ReadMessage(r.Body)
		if err != nil {
			log.Fatal(err)
		}

		
		from := m.Header.Get("From")
		from = strings.Split(from, "<")[1]
		from = strings.Split(from, ">")[0]
		ctx.Infof("Email replied from: %s", from)
		
		//***** Check if they should be responding or if someone is being snarky. ************
		groups, err := GetGroups(ctx)
		if err != nil {
			log.Fatal(err)
		}

		var g Group
		for _, group := range groups {
			for _, host := range group.Hosts {
				if strings.Contains(host.Emails, from) {
					g = group
					break;
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

		if !strings.Contains(g.Next.Emails, from) {
			log.Fatal("Sent from the wrong person!\n  Sent from: %s, but expected: %s", from, g.Hosts[0].Emails)
		}



        responseRegex := regexp.MustCompile(`(yes\b|no\b|skip\b)(.*?)`)
        body, err := ioutil.ReadAll(m.Body)
        ctx.Infof("email body: %s", body)

        if err != nil{
			log.Fatal(err)
		}
        

        s := string(body)
        bodyString := strings.ToLower(s)
        bodyString = strings.Split(bodyString, "it is your turn to host!")[0]

		switch responseRegex.FindString(bodyString) {
		case "yes":
			// Update the order in the group
        	hosts := g.Hosts
        	currentHost := hosts[0]
        	hosts = hosts[1:]
        	hosts = append(hosts, currentHost)  //Think slices are by reference??
        	g.Next = hosts[0]
        	g.save(ctx)
        	ctx.Infof("Match Yes")
		case "no":
			//Send an email to the next in line
	    	hosts := g.Hosts
	    	currentIndex := SliceIndex(len(hosts), func(i int) bool { return strings.Contains(hosts[i].Emails, from) }) 
	    	if(currentIndex < (len(hosts) - 1)) {
	    		g.Next = hosts[currentIndex + 1]
	    	} else {
	    		g.Next = hosts[0]
	    	}

	    	g.save(ctx)
	    	sendReminder(g, r)
	    	ctx.Infof("Match No")
		case "skip":
			//Respond with the current turn order for next week
	    	sendSkipMessage(g, r)
	        ctx.Infof("Match Skip")
		default:
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