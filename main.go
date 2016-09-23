package main

import (
	"net/http"
	"html/template"
	"appengine"
	"appengine/datastore"
	"appengine/user"
	"appengine/memcache"
	"time"
	"log"
	"bytes"
)


//<editor-fold defaultstate="collapsed"  desc="== datastore ==" >
func ticketKey(c appengine.Context) *datastore.Key {
	// The string "default_ticket" here could be varied to have multiple tickets.
	return datastore.NewKey(c, "Ticket", "default_ticket", 0, nil)
}

//</editor-fold>

//<editor-fold defaultstate="collapsed"  desc="== main view template ==" >
var mainPage = template.Must(template.New("guestbook").Parse(
	`<html><body><h1>Eddys Ticketing System</h1><h2>Tickets</h2>

		{{range .}}
			user: {{with .Author}}<b>{{.}}</b>{{else}}An anonymous person{{end}} <br/>
			time  <em>{{.Date.Format "3:04pm, Mon 2 Jan"}}</em><br/>
			content <blockquote>{{.Content}}</blockquote><br/>
		{{end}}

	<form action="/create" method="post">
	<div><textarea name="content" rows="3" cols="60"></textarea></div>
	<div><input type="submit" value="Create Ticket"></div>
	</form></body></html>
`))
//</editor-fold>

//<editor-fold defaultstate="collapsed"  desc="== schema ==" >
type Ticket struct {
	Author  string
	Content string
	Date    time.Time
}
//</editor-fold>

//<editor-fold defaultstate="collapsed"  desc="== init ==" >
func init() {
	http.HandleFunc("/", handleMainPage)
	http.HandleFunc("/create", handleTicket)

	http.HandleFunc("/tempstore", IndexHandler)
	http.HandleFunc("/mem", MemeCacheHandler)

}
//</editor-fold>

//<editor-fold defaultstate="collapsed"  desc="== handleMainPage  ==" >
func handleMainPage(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "GET requests only", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	c := appengine.NewContext(r)
	q := datastore.NewQuery("Ticket").Ancestor(ticketKey(c))

	var gg []*Ticket

	if _, err := q.GetAll(c, &gg); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if err := mainPage.Execute(w, gg); err != nil {
		c.Errorf("%v", err)
	}

}
//</editor-fold>

//<editor-fold defaultstate="collapsed"  desc="== datastore ==" >
func handleTicket(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "POST requests only", http.StatusMethodNotAllowed)
		return
	}
	c := appengine.NewContext(r)
	g := &Ticket{
		Content: r.FormValue("content"),
		Date:    time.Now(),
	}
	if u := user.Current(c); u != nil {
		g.Author = u.String()
	}
	key := datastore.NewIncompleteKey(c, "Ticket", ticketKey(c))
	if _, err := datastore.Put(c, key, g); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)

}
//</editor-fold>

//<editor-fold defaultstate="collapsed"  desc="== html form for memcache code==" >
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/index.html", http.StatusFound)
}
//</editor-fold>

//<editor-fold defaultstate="collapsed"  desc="== memcache ==" >
func MemeCacheHandler(w http.ResponseWriter, r *http.Request) {

	ctx := appengine.NewContext(r)
	var buf bytes.Buffer
	logger := log.New(&buf, "logger: ", log.Lshortfile)
	item := &memcache.Item{
		Key:   r.FormValue("email"),
		Value: []byte(r.FormValue("content")),
	}
	// Add the item to the memcache, if the key does not already exist
	if err := memcache.Add(ctx, item); err == memcache.ErrNotStored {
		logger.Print(ctx, "item with key %q already exists", item.Key)
	} else if err != nil {
		logger.Print(ctx, "error adding item: %v", err)
	}

}
//</editor-fold>
