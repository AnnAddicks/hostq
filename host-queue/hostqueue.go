package hostqueue

import (
  "appengine"
  "appengine/datastore"

  "github.com/GoogleCloudPlatform/go-endpoints/endpoints"

  "log"
  "net/http"
  "time"
  )


func init() {
  http.HandleFunc("/action/email/", SendEmail)
  http.HandleFunc("/api/group", CreateGroup)
}

type Host struct {
	Id int64 `json:"id" datastore:"-"`
	HostName string `json:"hostName"`
	Emails []string `json:"emails"`
	TimesHosted int64 `json:"timesHosted"`
	LastHosted time.Time `json:"lastHosted"`
}

type Group struct {
    Id   int64  `json: "id" datastore:"-"`
    GroupName string `json:"groupName"`
    GroupEmail string `json:"groupEmail"`
    Hosts []Host `json:"hosts"`
    Next Host `json:"next"`
}

//Datastore methods from:  http://stevenlu.com/posts/2015/03/23/google-datastore-with-golang/
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

func CreateGroup (w http.ResponseWriter, r *http.Request) {
  
}


