package hostqueue

import (
	"encoding/json"
	"net/http"
	"regexp"

	"appengine"

	"github.com/gorilla/mux"
)

type Status struct {
	Name  string `json:"name"`
	Hosts []Host `json:"hosts"`
	Next  int    `json:"next"`
}

//group/status/{groupUUID}
//Respond with the group name, current queue, who is up next in that queue, and if this week has a host yet.
func DisplayGroupStatus(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	pathVars := mux.Vars(r)

	uuid := pathVars["uuid"]
	c.Infof("UUID: %s", uuid)

	if isValidUUID((uuid)) {
		group, err := GetGroupByUUID(c, uuid)
		c.Infof("group: %v", group)
		if err != nil {
			panic(err)
		}

		status := convertToStatus(group)
		jsonStatus, err := json.Marshal(status)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonStatus)
	} else {
		w.Write([]byte("Invalid group"))
	}
}

func convertToStatus(group Group) Status {
	var status Status

	status.Name = group.GroupName
	status.Next = group.Next
	status.Hosts = group.Hosts

	return status
}

func isValidUUID(text string) bool {
	r := regexp.MustCompile("^[a-z0-9]{8}-[a-z0-9]{4}-[1-5][a-z0-9]{3}-[a-z0-9]{4}-[a-z0-9]{12}$")
	return r.MatchString(text)
}
