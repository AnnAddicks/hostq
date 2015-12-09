package hostqueue

import (
  "google.golang.org/appengine/datastore"
  "github.com/GoogleCloudPlatform/go-endpoints/endpoints"
  "golang.org/x/net/context"
  "log"
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

// A QuotesAPI struct defines all the endpoints of the quotes API.
// It will have functions for CRUD like Add, List etc.
type GroupAPI struct {
}

// Add creates a new quote given the fields in AddRequest, stores it in the
// datastore, and returns it.
func (GroupAPI) Add(c context.Context, g *Group) (*Group, error) {
  // We set the same parent key on every Quote entity to ensure each Quote
  // is in the same entity group. Queries across the single entity group
  // will be consistent.
  k := g.key(c)


  //TODO: Oh my, trusing input from a user!!
  k, err := datastore.Put(c, k, g)
  if err != nil {
    return nil, err
  }
  g.Id = k.IntID()
  return g, nil
}

//Datastore methods from:  http://stevenlu.com/posts/2015/03/23/google-datastore-with-golang/
func (group *Group) key(c context.Context) *datastore.Key {
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

func (group *Group) save(c context.Context) error {
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

func GetGroups(c context.Context) ([]Group, error) {
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

func SendEmail(c context.Context, w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/plain; charset=utf-8")
  g, err := GetGroups(c)

  if err != nil {
  	log.Fatal("Error getting groups: ", err)
  }

  for _, element := range g {
  	sendReminder(element, r) 
  }
}

func CreateGroup (w http.ResponseWriter, r *http.Request) {
  
}

func init() {
  //http.HandleFunc("/action/email/", SendEmail)
  //http.HandleFunc("/api/group", CreateGroup)

  // register the quotes API with cloud endpoints.
  api, err := endpoints.RegisterService(GroupAPI{}, "groupService", "v1", "Group API", true)
  if err != nil {
    panic(err)
  }

  info := api.MethodByName("Email").Info()
  info.Name, info.HTTPMethod, info.Path = "emailGroup", "GET", "groupService/email"

  info = api.MethodByName("Add").Info()
  info.Name, info.HTTPMethod, info.Path = "addGroup", "POST", "groupService"

  //http.HandleFunc("/_ah/mail/", incomingMail)
  info = api.MethodByName("Parse").Info()
  info.Name, info.HTTPMethod, info.Path = "parseEmail", "GET", "/_ah/mail/"

  // start handling cloud endpoint requests.
  endpoints.HandleHTTP()
}
