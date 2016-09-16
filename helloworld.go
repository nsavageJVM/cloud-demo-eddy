package helloworld

import (
	"fmt"
	"net/http"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func init() {
	http.HandleFunc("/", handle)
}



func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<html><body>Eddys Go Cloud Service </body></html>")
}
