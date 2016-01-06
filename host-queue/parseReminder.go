/*
	Script to handle an email response to the reminder.
	Steps:
		1. Get Sender - is it one of the registered senders in queue and are they the hosting group?
		2. Check if they should be responding or if someone is being snarky.
		3. Look for Yes/No/Skip
			3a. Yes - update the order in the group
			3b. No - send an email to the next in line, update group
			3c. Skip - respond with the current turn order for next week
	Notes: Super procedural right now.  I need to clean up the code

*/
package hostqueue

import (
	"appengine"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/mail"
	"regexp"
	"strings"
)

func IncomingMail(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	defer r.Body.Close()

	m, err := mail.ReadMessage(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	from := m.Header.Get("From")
	//clean up responses formatted with a name first.  Ex: 'Ann Addicks <test@test.com>'
	from = strings.Split(from, "<")[1]
	from = strings.Split(from, ">")[0]
	ctx.Infof("from : " + from)
	g, err := findGroup(ctx, from)
	if err != nil {
		log.Fatal(err)
	} else if isValidResponder(from, g) {
		responseRegex := regexp.MustCompile(`(yes\b|no\b|skip\b)(.*?)`)
		body, err := ioutil.ReadAll(m.Body)

		if err != nil {
			log.Fatal(err)
		}

		s := string(body)
		bodyString := strings.ToLower(s)
		bodyString = strings.Split(bodyString, "it is your turn to host!")[0]

		switch responseRegex.FindString(bodyString) {
		case "yes":
			// Update the order in the group by deleting the current host and appending them to the end.
			hosts := g.Hosts
			currentHost := hosts[g.Next]
			currentHost.TimesHosted++

			hosts = append(hosts[:g.Next], hosts[g.Next+1:]...)
			hosts = append(hosts, currentHost)
			g.Hosts = hosts
			g.Next = 0

			g.save(ctx)
			sendHostConfirmedMessage(g, r)
			ctx.Infof("Match Yes")
		case "no":
			//Send an email to the next in line && update the group
			hosts := g.Hosts

			if g.Next < (len(hosts) - 1) {
				g.Next++
			} else {
				g.Next = 0
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
	} else {
		ctx.Infof("Sender is not valid")
	}
}

/*Pull all the groups and loop through to find the email address.
  I need to learn more about querying app engines datastore to make this nicer.*/
func findGroup(ctx appengine.Context, from string) (Group, error) {
	groups, err := GetGroups(ctx)
	var g Group
	if err != nil {
		return g, err
	}

	for _, group := range groups {
		for _, host := range group.Hosts {
			if strings.Contains(host.Emails, from) {
				return group, nil
			}
		}
	}
	return g, fmt.Errorf("Cannot find a group for: %s", from)
}

//Check if they should be responding or if someone is being snarky.
func isValidResponder(from string, g Group) bool {
	validEmails := g.Hosts[g.Next].Emails
	return strings.Contains(validEmails, from)
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
