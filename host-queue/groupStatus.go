package hostqueue

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"appengine"

	"github.com/gorilla/mux"
)

//group/status/{groupUUID}
//Respond with the group name, current queue, who is up next in that queue, and if this week has a host yet.
func DisplayGroupStatus(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	pathVars := mux.Vars(r)
	AddUuid(c)

	uuid := pathVars["uuid"]
	c.Infof("UUID: %s", uuid)

	if isValidUUID((uuid)) {
		group, err := GetGroupByUUID(c, uuid)
		c.Infof("group: %v", group)
		if err != nil {
			//c.Infof()
			w.Write([]byte("invalid group"))
		}

		jsonGroup, err := json.Marshal(group)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonGroup)
	}

	w.Write([]byte("invalid group"))
}

func isValidUUID(text string) bool {
	r := regexp.MustCompile("^[a-z0-9]{8}-[a-z0-9]{4}-[1-5][a-z0-9]{3}-[a-z0-9]{4}-[a-z0-9]{12}$")
	return r.MatchString(text)
}
