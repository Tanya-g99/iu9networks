package main

import (
	"bufio"
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"

	log2 "github.com/mgutz/logxi/v1"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log2.Info("connecting to ", "server", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log2.Error(err.Error())
		return
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			log2.Info("Enter the task")
			_, message, err := c.ReadMessage()
			if err != nil {
				log2.Error(err.Error())
				return
			}

			log2.Info("Answer got successful", "answer: ", string(message))
		}
	}()

	outgoing := make(chan string)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			line, _, err := reader.ReadLine()
			if err != nil {
				break
			}
			outgoing <- string(line)
		}
	}()

	for {

		select {
		case <-done:
			return
		case line := <-outgoing:
			err := c.WriteMessage(websocket.TextMessage, []byte(line))
			if err != nil {
				log2.Error(err.Error())
				return
			}
		case <-interrupt:
			log2.Info("Interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log2.Error(err.Error())
				return
			}
			select {
			case <-done:
			case <-outgoing:
			}
			return
		}
	}
}
