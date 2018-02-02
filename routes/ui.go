package routes

import (
	"html/template"
	"net/http"

	"github.com/repejota/misto"
)

// HadleHome ...
func HandleHome(w http.ResponseWriter, r *http.Request) {
	data, err := misto.Asset("assets/index.html")
	if err != nil {
		panic(err)
	}
	t, _ := template.New("home").Parse(string(data))
	t.Execute(w, nil)
}
