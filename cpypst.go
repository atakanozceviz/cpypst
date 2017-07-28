package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

var re = regexp.MustCompile(`:[0-9]+`)

func connect(w http.ResponseWriter, r *http.Request) {
	addr := r.RemoteAddr
	name := r.FormValue("name")
	if name != "" && addr != "" {
		ip := re.ReplaceAllString(addr, "")
		ui.Connections.Add(model.Connection{ip, name, true})
		log.Println(name + " (" + ip + ") is connected!")
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
	addr := ui.Connections.Connections

	if len(addr) > 0 {
		for _, v := range addr {
			if v.Active == true {
				req, err := http.NewRequest("POST", "http://"+v.Addr+":8080/paste", bytes.NewBuffer(clip))
				if err != nil {
					log.Println(err)
				}
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
	action := r.FormValue("action")
	if action == "dlt" {
		id, err := strconv.Atoi(r.FormValue("ID"))
		if err != nil {
			log.Println(err)
		}
		ui.History.Remove(id)
	}
	if action == "cpy" {
		id, err := strconv.Atoi(r.FormValue("ID"))
		if err != nil {
			log.Println(err)
		}
		cpy := ui.History.History[id].Content
		tmp.Write(cpy)
		clipboard.WriteAll(cpy)
	}
	if action == "enable" {
		id := r.FormValue("ID")
		if id != "" {
			ui.Connections.Connections[id].Active = true
			tmp.Write("")
		}
	}
	if action == "disable" {
		id := r.FormValue("ID")
		if id != "" {
			ui.Connections.Connections[id].Active = false
		}
	}
}

func main() {
	go func() {
		var name string
		var addr string

		fmt.Print("Enter your name: ")
		fmt.Scanln(&name)
		fmt.Print("Enter ip address: ")
		fmt.Scanln(&addr)
		_, err := http.Get("http://" + addr + ":8080/connect?name=" + name)
		if err != nil {
			log.Println(err)
		} else {
			go checkClip()
		}
	}()

	http.HandleFunc("/", ui.HistoryUI)
	http.HandleFunc("/connections", ui.ConnectionsUI)
	http.HandleFunc("/jquery", ui.Jquery)

	http.HandleFunc("/action", action)

	http.HandleFunc("/connect", connect)
	http.HandleFunc("/paste", paste)

	fmt.Println("Serving")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
