package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	. "lab3_2/proto"

	log "github.com/mgutz/logxi/v1"
)

var AllConnections map[string]bool

func handleRegister(w http.ResponseWriter, r *http.Request) {
	var conf Node

	conf.Addr = r.RemoteAddr
	for addr := range AllConnections {
		conf.Connections = append(conf.Connections, addr)
	}
	encoded, _ := json.Marshal(conf)
	AllConnections[conf.Addr] = true
	w.Write(encoded)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	delete(AllConnections, r.RemoteAddr)
}

func handleUpdatePeers(w http.ResponseWriter, r *http.Request) {
	var conf Node

	for addr := range AllConnections {
		conf.Connections = append(conf.Connections, addr)
	}
	encoded, err := json.Marshal(conf)
	if err != nil {
		log.Error(err.Error())
		return
	}
	w.Write(encoded)
}

var requestNum int

func LogRequest(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestNum++
		log.Info(fmt.Sprintf("Got new request #%d", requestNum), "address", r.RemoteAddr)
		f(w, r)
		log.Info(fmt.Sprintf("Request #%d processed succesfully", requestNum))
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/register", LogRequest(handleRegister))
	mux.HandleFunc("/logout", handleLogout)
	mux.HandleFunc("/updatepeers", handleUpdatePeers)

	AllConnections = make(map[string]bool)

	server := http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: mux,
	}
	log.Info("Starting server on 0.0.0.0:8000")
	err := server.ListenAndServe()
	if err != nil {
		log.Error(err.Error())
	}
}
