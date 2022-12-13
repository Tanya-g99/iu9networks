package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
	log1 "github.com/mgutz/logxi/v1"
)

var e = 0.001

var addr = flag.String("addr", "localhost:8082", "http service address")
var upgrader = websocket.Upgrader{} // use default options
func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		task := string(message)
		coefficients := strings.Split(task, " ")
		log1.Info("request: %s", task)

		if len(coefficients) == 3 {
			a, err := strconv.ParseFloat(coefficients[0], 64)
			if err != nil {
				log1.Error(err.Error())
				return
			}
			b, err := strconv.ParseFloat(coefficients[1], 64)
			if err != nil {
				log1.Error(err.Error())
				return
			}
			c, err := strconv.ParseFloat(coefficients[2], 64)
			if err != nil {
				log1.Error(err.Error())
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
		fmt.Printf("Answer: %s\n", task)
		err = c.WriteMessage(mt, []byte(task))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}
func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	log.Println(http.ListenAndServe(*addr, nil))
	log.Fatal(http.ListenAndServe(*addr, nil))
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<link
	href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/css/bootstrap.min.css"
    rel="stylesheet"
    integrity="sha384-EVSTQN3/azprG1Anm3QDgpJLIm9Nao0Yz1ztcQTwFspd3yD65VohhpuuCOmLASjC"
	crossorigin="anonymous"
/>
<script>
window.addEventListener("load", function(evt) {
var output = document.getElementById("output");
var input = document.getElementById("input");
var ws;
var print = function(message) {
var d = document.createElement("div");
d.textContent = message;
output.appendChild(d);
output.scroll(0, output.scrollHeight);
};
if (ws) {
return false;
}
ws = new WebSocket("{{.}}");
ws.onopen = function(evt) {
document.getElementById("send").onclick = function(evt) {
if (!ws) {
return false;
}
print("SEND: " + input.value);
ws.send(input.value);
return false;
};
}
ws.onclose = function(evt) {
print("CLOSE CONSOLE");
ws = null;
}
ws.onmessage = function(evt) {
print("RESPONSE: " + evt.data);
}
ws.onerror = function(evt) {
print("ERROR: " + evt.data);
}
return false;
document.getElementById("close").onclick = function(evt) {
if (!ws) {
return false;
}
ws.close();
return false;
};
});
</script>
</head>
<body>
<form class="text-center m-3">
<p class=""><input id="input" type="text" value="1 2 1">
<p>
<button class="btn btn-secondary mx-auto" id="send">Отправить</button>
</form>
<div class="text-center m-3 pt-3" id="output" style="max-height: 70vh; overflow-y: scroll;"></div>
</body>
</html>
`))
