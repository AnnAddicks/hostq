package hostqueue

import (
  "appengine"
  "appengine/datastore"
  "encoding/json"
  "net/http"
  "time"
  )


type Host struct {
	HostName string `json:"hostName"`
	Emails string `json:"emails"`
	TimesHosted int64 `json:"timesHosted"`
	LastHosted time.Time `json:"lastHosted"`
}

type Group struct {
    Id   int64  `json: "id" datastore:"-"`
    GroupName string `json:"groupName"`
    GroupEmail string `json:"groupEmail"`
    Hosts []Host `json:"hosts"`
    Next int `json:"next"`
}



// Add creates a new group, stores it in the
// datastore, and returns it.
func  Add(w http.ResponseWriter, r *http.Request)  {
  // We set the same parent key on every Group entity to ensure each Group
  // is in the same entity group. Queries across the single entity group
  // will be consistent.
  if r.Method == "POST" {
    ctx := appengine.NewContext(r)
    var group *Group

    err := json.NewDecoder(r.Body).Decode(&group)
    if err != nil {
      panic(err)
    }
    k := group.key(ctx)


    //TODO: Oh my, trusing input from a user!!
    k, err = datastore.Put(ctx, k, group)
    if err != nil {
      w.Header().Set("Content-Type", "application/json")
      panic(err)
    }
    group.Id = k.IntID()
    g, err := json.Marshal(group)
    if err != nil {
      panic(err)  
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(g)
  } else {
     http.Error(w, "Invalid request method.", 405)
  }
}


func (group *Group) key(c appengine.Context) *datastore.Key {
  return datastore.NewKey(c, "Group", "default_group", 0, nil)
}

func (group *Group) save(ctx appengine.Context) error {
  // reference the key function and generate it
  // accordingly basically its isNew true/false
  k, err := datastore.Put(ctx, group.key(ctx), group)
  if err != nil {
    return err
  }

  // The Id on the model is not prepopulated so we'll have
  // to append manually
  group.Id = k.IntID()
  return nil
}

func GetGroups(ctx appengine.Context) ([]Group, error) {
  q := datastore.NewQuery("Group")

  var groups []Group
  keys, err := q.GetAll(ctx, &groups)
  if err != nil {
    return nil, err
  }

  
  for i := 0; i < len(groups); i++ {
    groups[i].Id = keys[i].IntID()
  }

  return groups, nil
}

func  SendEmail(w http.ResponseWriter, r *http.Request) {
  ctx := appengine.NewContext(r)
  w.Header().Set("Content-Type", "text/plain; charset=utf-8")
  g, err := GetGroups(ctx)

  if err != nil {
  	panic(err)
  }

  ctx.Infof("Groups: %v", g)
  for _, element := range g {
  	sendReminder(element, r) 
  }
}


func init() {
   http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request)  {
      w.Write([]byte("hostq")) 
   })

   http.HandleFunc("/_ah/mail/reminder@hostqueue-1146.appspotmail.com", IncomingMail)
   http.HandleFunc("/group/add", Add)
   http.HandleFunc("/group/action/email", SendEmail)
}
