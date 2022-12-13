package proto

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	log "github.com/mgutz/logxi/v1"
)

const url string = "http://localhost:8000"

type Integral struct {
	Expr string  `json:"expr"`
	A    float64 `json:"a"`
	B    float64 `json:"b"`
}

func (i *Integral) Calculate() float64 {
	expr := strings.Split(i.Expr, " ")
	if len(expr) < 2 {
		if expr[0] == "x" {
			return i.B*i.B - i.A*i.A
		}
		if len(i.Expr) > 5 {
			switch i.Expr[:4] {
			case "sin(":
				return math.Cos(i.B) - math.Cos(i.A)
			case "cos(":
				return math.Sin(i.B) - math.Sin(i.A)
			}
		}
		x, err := strconv.Atoi(expr[0])
		if err != nil {
			fmt.Println("unknown integral, please try again")
			return 0
		}
		return float64(x) * (i.B - i.A)
	}
	x1 := Integral{Expr: expr[0], A: i.A, B: i.B}
	x2 := Integral{Expr: expr[2], A: i.A, B: i.B}
	switch expr[1] {
	case "+":
		return x1.Calculate() + x2.Calculate()
	case "-":
		return x1.Calculate() - x2.Calculate()
	case "/":
		{
			if expr[0] != "1" {
				fmt.Println("unknown integral, please try again")
				return 0
			}
			return math.Log(i.B) - math.Log(i.A)
		}
	case "*":
		{
			y1, err1 := strconv.Atoi(expr[0])
			y2, err2 := strconv.Atoi(expr[2])
			if err1 == nil {
				if err2 == nil {
					return float64(y1*y2) * (i.B - i.A)
				}
				x := Integral{Expr: expr[2], A: i.A, B: i.B}
				return float64(y1) * x.Calculate()
			}
			if err2 == nil {
				x := Integral{Expr: expr[0], A: i.A, B: i.B}
				return float64(y2) * x.Calculate()
			}
			return (i.B*i.B*i.B - i.A*i.A*i.A) / 3
		}
	}

	fmt.Println("unknown integral, please try again")
	return 0
}

// Config - параметры пира
type Node struct {
	Addr        string   `json:"addr"`
	Connections []string `json:"connections"`
}

type Package struct {
	To   string   `json:"to"`
	From string   `json:"from"`
	Data Integral `json:"integral"`
}

func (n *Node) UpdatePeers() {
	resp, err := http.Get(url + "/updatepeers")
	if err != nil {
		log.Error("request failed", "reason", err.Error())
	} else {
		body, _ := bufio.NewReader(resp.Body).ReadBytes('\u00C6')
		var foo Node
		if err := json.Unmarshal(body, &foo); err != nil {
			log.Error(err.Error())
			return
		}
		n.Connections = nil
		for _, v := range foo.Connections {
			if v != n.Addr {
				n.Connections = append(n.Connections, v)
			}
		}
	}
}

func (n *Node) Run(handleServer func(*Node), handleClient func(*Node)) {
	go handleServer(n)
	handleClient(n)
}

func (n *Node) SendToAll(message Integral) {
	var new_pack = Package{
		From: n.Addr,
		Data: message,
	}
	count := len(n.Connections)
	ab := (message.B - message.A) / float64(count)
	new_pack.Data.A = message.A
	new_pack.Data.B = new_pack.Data.A
	for _, addr := range n.Connections {
		new_pack.Data.A = new_pack.Data.B
		new_pack.Data.B += ab
		new_pack.To = addr
		n.Send(new_pack)
	}
}

func (n *Node) Send(pack Package) {
	json_packet, err := json.Marshal(pack)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = http.Post("http://"+pack.To[:len(pack.To)-2]+string(pack.To[len(pack.To)-1]), "application/json", strings.NewReader(string(json_packet)))
	if err != nil {
		fmt.Println(err)
	}
}
