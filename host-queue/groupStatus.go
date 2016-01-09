package hostqueue

import (
	"html/template"
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

/* group/status/{groupUUID}
   Respond with the group name, current queue, who is up next in that queue, and if this week has a host yet.
   I'm not sure what this app is going to be yet, api with spa? or server & client code
*/
func DisplayGroupStatus(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	pathVars := mux.Vars(r)

	uuid := pathVars["uuid"]
	ctx.Infof("UUID: %s", uuid)

	if isValidUUID((uuid)) {
		group, err := GetGroupByUUID(ctx, uuid)
		ctx.Infof("group: %v", group)
		if err != nil {
			panic(err)
		}

		status := convertToStatus(group)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		var tpl = template.Must(template.ParseGlob("templates/*.html"))
		if err := tpl.ExecuteTemplate(w, "index.html", status); err != nil {
			ctx.Infof("%v", err)
		}

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
