package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	log2 "github.com/mgutz/logxi/v1"

	"github.com/gorilla/websocket"
)

var e = 0.001
var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log2.Error(err.Error())
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		task := string(message)
		if err != nil {
			log2.Error(err.Error())
			break
		}
		log2.Info("Task got successful", "task: ", task)
		coefficients := strings.Split(task, " ")

		if len(coefficients) == 3 {
			a, err := strconv.ParseFloat(coefficients[0], 64)
			if err != nil {
				log2.Error(err.Error())
				return
			}
			b, err := strconv.ParseFloat(coefficients[1], 64)
			if err != nil {
				log2.Error(err.Error())
				return
			}
			c, err := strconv.ParseFloat(coefficients[2], 64)
			if err != nil {
				log2.Error(err.Error())
				return
			}

			D := b*b - 4*a*c
			if D > -e && D < e {
				res := -b / (2 * a)
				task = fmt.Sprintf("Solution is %f", res)
			} else if D < 0 {
				task = "Not solution"
			} else {
				res1, res2 := (-b+math.Sqrt(D))/(2*a), (-b-math.Sqrt(D))/(2*a)
				task = fmt.Sprintf("Solution is %f and %f", res1, res2)
			}
		} else {
			task = "Entered data's flawed"
		}
		err = c.WriteMessage(mt, []byte(task))
		if err != nil {
			log2.Error(err.Error())
			break
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	log2.Info(http.ListenAndServe(*addr, nil).Error())
}
