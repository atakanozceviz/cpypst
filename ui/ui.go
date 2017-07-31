package ui

import (
	//"html/template"
	"log"
	"net/http"

	"github.com/arschles/go-bindata-html-template"
	"github.com/atakanozceviz/cpypst/model"
)

var History model.History
var Incoming model.Connections
var Outgoing model.Connections

func HistoryUI(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.New("ui", Asset).Parse("ui/history.html"))
	//tpl := template.Must(template.ParseGlob("./ui/*"))

	pData := struct {
		History  []model.HistItem
		Incoming map[string]*model.Connection
	}{
		History:  History.History,
		Incoming: Incoming.Connections,
	}

	tpl.ExecuteTemplate(w, "ui", pData)
}

func ConnectionsUI(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.New("ui", Asset).Parse("ui/connections.html"))
	//tpl := template.Must(template.ParseGlob("./ui/*"))
	Connections := struct {
		Incoming map[string]*model.Connection
		Outgoing map[string]*model.Connection
	}{
		Incoming.Connections,
		Outgoing.Connections,
	}

	tpl.ExecuteTemplate(w, "ui", Connections)
}

func Jquery(w http.ResponseWriter, r *http.Request) {
	jquery, err := Asset("ui/jquery.js")
	if err != nil {
		log.Println(err)
	}
	w.Write(jquery)

}
