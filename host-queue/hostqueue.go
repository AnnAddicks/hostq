package hostqueue

import (
  "appengine"
  "appengine/datastore"
  "encoding/json"
  "github.com/go-martini/martini"
  "net/http"
  "time"
  )


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



// Add creates a new quote given the fields in AddRequest, stores it in the
// datastore, and returns it.
func  Add(c martini.Context, r *http.Request) (*Group, error) {
  // We set the same parent key on every Quote entity to ensure each Quote
  // is in the same entity group. Queries across the single entity group
  // will be consistent.
  ctx := appengine.NewContext(r)

  var group *Group

  err := json.NewDecoder(r.Body).Decode(&group)
  if err != nil {
    panic(err)
  }
  k := g.key(ctx)


  //TODO: Oh my, trusing input from a user!!
  k, err := datastore.Put(ctx, k, g)
  if err != nil {
    return nil, err
  }
  g.Id = k.IntID()
  return g, nil
}

//Datastore methods from:  http://stevenlu.com/posts/2015/03/23/google-datastore-with-golang/
func (group *Group) key(c appengine.Context) *datastore.Key {
  // if there is no Id, we want to generate an "incomplete"
  // one and let datastore determine the key/Id for us
  //if group.Id == 0 {
  //  return datastore.NewIncompleteKey(c, "Group", nil)
  //}

  // if Id is already set, we'll just build the Key based
  // on the one provided.
  //return datastore.NewKey(c, "Group", "", group.Id, nil)

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

func  SendEmail(c martini.Context, w http.ResponseWriter, r *http.Request) {
  ctx := appengine.NewContext(r)
  w.Header().Set("Content-Type", "text/plain; charset=utf-8")
  g, err := GetGroups(ctx)

  if err != nil {
  	panic("Error getting groups: ", err)
  }

  for _, element := range g {
  	sendReminder(element, r) 
  }
}


func init() {
   m := martini.Classic()
   
   m.Get("/_ah/mail/", IncomingMail)
   m.Post("/group/add", Add)
   m.Get("/group/action/email", SendEmail)

   m.Get("/", func() string {
     return "hostq"
   })


   http.Handle("/", m)
}
