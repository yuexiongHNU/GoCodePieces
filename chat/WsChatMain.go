package main

import (
	"flag"
	"net/http"
	"log"
)

var addr = flag.String("addr", ":8080", "http service address")

func serverHome(w http.ResponseWriter, r *http.Request)  {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w,"Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	log.Println("here")
	http.ServeFile(w, r, "C:\\Users\\Administrator\\go\\src\\HelloWorld\\chat\\home.html")
}

func main()  {
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", serverHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serverWs(hub, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("Listen Server error:", err)
	}
}