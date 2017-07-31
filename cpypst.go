package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/atakanozceviz/cpypst/model"
	"github.com/atakanozceviz/cpypst/ui"
	"github.com/atotto/clipboard"
)

const (
	port = "8080"
)

var tmp model.Tmp
var clip model.Tmp
var lname, _ = os.Hostname()

var re = regexp.MustCompile(`:[0-9]+`)

func connect(w http.ResponseWriter, r *http.Request) {
	addr := r.RemoteAddr
	name := r.FormValue("name")

	if name != "" && addr != "" {
		ip := re.ReplaceAllString(addr, "")
		ui.Incoming.Add(model.Connection{ip, name, true})
		w.Write([]byte(lname))
		log.Println(name + " (" + ip + ") is connected!")
	} else {
		io.WriteString(w, "Wrong request!")
	}
}

func checkClip() {
	var err error
	var x string
	for {
		x, err = clipboard.ReadAll()
		if err != nil {
			log.Println(err)
		}
		clip.Write(x)
		if clip.Read() != tmp.Read() {
			tmp.Write(clip.Read())
			send([]byte(clip.Read()))
		}
		time.Sleep(time.Second * 1)
	}
}

func send(clip []byte) {
	client := &http.Client{}
	addr := ui.Outgoing.Connections

	if len(addr) > 0 {
		for _, v := range addr {
			if v.Active == true {
				req, err := http.NewRequest("POST", "http://"+v.Ip+":8080/paste", bytes.NewBuffer(clip))
				if err != nil {
					log.Println(err)
				}
				req.Header.Set("Send", "true")
				resp, err := client.Do(req)
				if err != nil {
					log.Println(err)
				}
				resp.Body.Close()
			}
		}
	}
}

func paste(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	rbody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	clip.Write(string(rbody))
	tmp.Write(string(rbody))
	if clipboard.WriteAll(clip.Read()) != nil {
		log.Println(err)
	}
	ip := re.ReplaceAllString(r.RemoteAddr, "")
	ui.History.Add(model.HistItem{Ip: ip, Content: clip.Read()})

}

func action(w http.ResponseWriter, r *http.Request) {
	switch action := r.FormValue("action"); action {
	case "dlt":
		{
			id, err := strconv.Atoi(r.FormValue("ID"))
			if err != nil {
				log.Println(err)
			} else {
				ui.History.Remove(id)
			}
		}
	case "cpy":
		{
			id, err := strconv.Atoi(r.FormValue("ID"))
			if err != nil {
				log.Println(err)
			} else {
				cpy := ui.History.History[id].Content
				tmp.Write(cpy)
				clipboard.WriteAll(cpy)
			}
		}
	case "ienable":
		{
			id := r.FormValue("ID")
			if id != "" {
				if val, ok := ui.Incoming.Connections[id]; ok {
					val.Active = true
				} else {
					io.WriteString(w, "Couldn't find connection")
				}
			}
		}
	case "idisable":
		{
			id := r.FormValue("ID")
			if id != "" {
				if val, ok := ui.Incoming.Connections[id]; ok {
					val.Active = false
				} else {
					io.WriteString(w, "Couldn't find connection")
				}
			}
		}
	case "oenable":
		{
			id := r.FormValue("ID")
			if id != "" {
				if val, ok := ui.Outgoing.Connections[id]; ok {
					val.Active = true
				} else {
					io.WriteString(w, "Couldn't find connection")
				}
			}
		}
	case "odisable":
		{
			id := r.FormValue("ID")
			if id != "" {
				if val, ok := ui.Outgoing.Connections[id]; ok {
					val.Active = false
				} else {
					io.WriteString(w, "Couldn't find connection")
				}
			}
		}
	default:
		io.WriteString(w, "Invalid action!")
	}

}

func main() {
	go func() {
		var addr string
		for {
			fmt.Print("Enter ip address to add a connection: ")
			fmt.Scanln(&addr)
			req, err := http.Get("http://" + addr + ":8080/connect?name=" + lname)
			req.Header.Set("Connect", "true")
			if err != nil {
				log.Println(err)
			} else {
				body, err := ioutil.ReadAll(req.Body)
				if err != nil {
					log.Println(err)
					req.Body.Close()
				} else {
					ui.Outgoing.Add(model.Connection{addr, string(body), true})
					req.Body.Close()
					go checkClip()
				}
			}
		}
	}()

	http.HandleFunc("/", ui.HistoryUI)
	http.HandleFunc("/connections", ui.ConnectionsUI)
	http.HandleFunc("/jquery", ui.Jquery)

	http.HandleFunc("/action", action)

	http.HandleFunc("/connect", connect)
	http.HandleFunc("/paste", paste)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
