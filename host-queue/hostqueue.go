package hostqueue

import (
  "appengine"
  "appengine/datastore"

  "net/http"
  "time"
  )


func init() {
  http.HandleFunc("/action/email/", SendEmail)
}

type Host struct {
	Id int64 `datastore:"-"`
	HostName string
	Email [] string
	TimesHosted int64
	LastHosted time.Time
}

type Group struct {
    Id   int64  `datastore:"-"`
    GroupName string
    GroupEmail string
    Hosts [] Host
}

func GetGroups(c appengine.Context) ([]Group, error) {
  q := datastore.NewQuery("Group")
  
  var groups []Group
  keys, err := q.GetAll(c, &groups)
  if err != nil {
    return nil, err
  }

  
  for i := 0; i < len(groups); i++ {
    groups[i].Id = keys[i].IntID()
  }

  return groups, nil
}

func SendEmail(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  w.Header().Set("Content-Type", "text/plain; charset=utf-8")

  c.Infof("Getting all groups")


  c.Infof("Sending Emails")



}