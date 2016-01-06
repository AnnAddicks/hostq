package hostqueue

import (
	"appengine"
	"net/http"
)

//group/status/{groupUUID}
//Respond with the group name, current queue, who is up next in that queue, and if this week has a host yet.
func DisplayGroupStatus(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	c.Infof("URI", r.URL.Path[1:])
}
