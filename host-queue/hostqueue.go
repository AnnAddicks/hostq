package hostqueue

import (
  "appengine"
  "appengine/datastore"

  "log"
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

func (group *Group) key(c appengine.Context) *datastore.Key {
  // if there is no Id, we want to generate an "incomplete"
  // one and let datastore determine the key/Id for us
  if group.Id == 0 {
    return datastore.NewIncompleteKey(c, "Group", nil)
  }

  // if Id is already set, we'll just build the Key based
  // on the one provided.
  return datastore.NewKey(c, "Group", "", group.Id, nil)
}

func (group *Group) save(c appengine.Context) error {
  // reference the key function and generate it
  // accordingly basically its isNew true/false
  k, err := datastore.Put(c, group.key(c), group)
  if err != nil {
    return err
  }

  // The Id on the model is not prepopulated so we'll have
  // to append manually
  group.Id = k.IntID()
  return nil
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
  g, err := GetGroups(c)

  if err != nil {
  	log.Fatal("Error getting groups: ", err)
  }

  c.Infof("Sending Emails")
  for _, element := range g {
  	sendReminder(element, r) 
  }
}


