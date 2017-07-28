package ui

import (
	"log"
	"net/http"

	"github.com/arschles/go-bindata-html-template"
	"github.com/atakanozceviz/cpypst/model"
)

var History model.History
var Connections model.Connections

func HistoryUI(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.New("ui", Asset).Parse("ui/history.html"))
	//tpl := template.Must(template.ParseGlob("./ui/*"))

	tpl.ExecuteTemplate(w, "ui", History)
}

func ConnectionsUI(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.New("ui", Asset).Parse("ui/connections.html"))
	//tpl := template.Must(template.ParseGlob("./ui/*"))

	tpl.ExecuteTemplate(w, "ui", Connections)
}

func Jquery(w http.ResponseWriter, r *http.Request) {
	jquery, err := Asset("ui/jquery.js")
	if err != nil {
		log.Println(err)
	}
	w.Write(jquery)

}
