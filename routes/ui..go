package routes

import (
	"html/template"
	"net/http"
)

// HadleHome ...
func HandleHome(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("index.html")
	t.Execute(w, nil)
}
