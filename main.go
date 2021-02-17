package main

import (
	"flag"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "public/index.html")
}

func catchData(hub *Hub, w http.ResponseWriter, r *http.Request) {
	//log.Println("Data:", r.URL)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Inc r.Body read: ", err)
	} else {
		hub.broadcast <- []byte(body)
	}
	http.ServeFile(w, r, "data/thanks.json")


}
func catchUpdate(hub *Hub, w http.ResponseWriter, r *http.Request) {
	//log.Println("Update:" , r.URL)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal("r.Body read: ", err)
	}
	hub.cached = body
	hub.broadcast <- []byte(body)

	http.ServeFile(w, r, "data/thanks.json")
}


func runSocks(){
	flag.Parse()
	hub := newHub()

	dat,_ := ioutil.ReadFile("data/initial.json")
	hub.cached = dat

	go hub.run()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	router.HandleFunc("/inc", func(w http.ResponseWriter, r *http.Request) {
		catchData(hub, w, r)
	})
	router.HandleFunc("/upd", func(w http.ResponseWriter, r *http.Request) {
		catchUpdate(hub, w, r)
	})
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

	log.Fatal(http.ListenAndServe(":8771", router))
}

func main() {
	runSocks()
}