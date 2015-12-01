package hostqueue

func init() {
        http.HandleFunc("/_ah/mail/", incomingMail)
}

func incomingMail(w http.ResponseWriter, r *http.Request) {
        //Sample from https://cloud.google.com/appengine/docs/go/mail/
        c := appengine.NewContext(r)
        defer r.Body.Close()
        var b bytes.Buffer
        if _, err := b.ReadFrom(r.Body); err != nil {
                c.Errorf("Error reading body: %v", err)
                return
        }
        c.Infof("Received mail: %v", b)

        //Steps
        //1. Get Sender - is it one of the registered senders in queue and are they the hosting group?
        //2. Look for Yes/No/Skip
        //3. If yes or skip, respond with the current turn order
        //4. If No send an email to the next in line
}