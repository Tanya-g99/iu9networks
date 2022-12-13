package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	. "lab3_2/proto"

	log "github.com/mgutz/logxi/v1"
)

const url string = "http://localhost:8000"

func main() {
	var conf Node
	resp, err := http.Get(url + "/register")
	if err != nil {
		log.Error(err.Error())
		return
	}
	body, _ := bufio.NewReader(resp.Body).ReadBytes('\u00C6')
	if err := json.Unmarshal(body, &conf); err != nil {
		log.Error(err.Error())
		return
	}

	conf.Run(handleServer, handleClient)
}

func handleServer(n *Node) {
	mux := http.NewServeMux()

	handleConnection := func(w http.ResponseWriter, r *http.Request) {
		var pack Package
		body, err := bufio.NewReader(r.Body).ReadBytes('\u00C6')
		if err != nil && err != io.EOF {
			log.Error(err.Error())
			fmt.Fprintf(w, "I'm chereshnya")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := json.Unmarshal(body, &pack); err != nil {
			log.Error(err.Error())
			fmt.Fprintf(w, "I'm teepot")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Info("task from ", pack.From)
		log.Info("desc", "integral ", pack.Data.Expr, "upper bound ", pack.Data.A, "lower bound", pack.Data.B, "result", pack.Data.Calculate())
	}

	mux.HandleFunc("/", handleConnection)
	server := http.Server{
		Addr:    n.Addr[:len(n.Addr)-2] + string(n.Addr[len(n.Addr)-1]),
		Handler: mux,
	}
	go server.ListenAndServe()
}

func handleClient(n *Node) {
	for {
		message := InputString()
		splited := strings.Split(message, " ")
		switch splited[0] {
		case "exit":
			http.Get(url + "/logout")
			os.Exit(0)
		case "calc":
			n.UpdatePeers()
			fmt.Print("enter integral: ")
			var a, b float64
			integral, _ := bufio.NewReader(os.Stdin).ReadString('\n')
			fmt.Print("enter bounds: ")
			fmt.Scan(&a, &b)
			n.SendToAll(Integral{
				Expr: integral,
				A:    a,
				B:    b,
			})

		default:
			log.Warn("unknown command, please try again")
		}
	}
}

func InputString() string {
	msg, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.Replace(msg, "\n", "", -1)
}
