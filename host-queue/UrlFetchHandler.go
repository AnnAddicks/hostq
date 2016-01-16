package hostqueue

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
	"net/http"
)

func handleURLFetch(req *http.Request, r *http.Request) {

	ctx := appengine.NewContext(r)
	client := urlfetch.Client(ctx)

	_, err := client.Do(req)
	if err != nil {
		panic(err)
	}
}
