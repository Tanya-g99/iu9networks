package main

import (
	"bytes"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gorilla/websocket"
	log1 "github.com/mgutz/logxi/v1"
)

var hashTable = map[string]string{
	"tupic":      "1234",
	"root":       "12",
	"first_girl": "gosha",
	"aboba":      "03",
}

func handle(data []string) string {

	name := data[0]
	args := data[1:]

	log1.Info("New command", "name", name, "args", args)

	var buf bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &buf

	err := cmd.Run()
	if err != nil {
		log1.Error(err.Error())
	}

	return buf.String()
}

var addr = flag.String("addr", "localhost:8083", "http service address")
var upgrader = websocket.Upgrader{} // use default options
func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log1.Error(err.Error())
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log1.Info("request:", string(message))
		msg := strings.Split(string(message), "~")
		if hashTable[msg[1]] == msg[2] {
			answer := handle(strings.Split(msg[0], " "))

			if answer != "" {
				err = c.WriteMessage(mt, []byte(answer))
			}
			if err != nil {
				log1.Error(err.Error())
				break
			}
		} else {
			err = c.WriteMessage(mt, []byte("User not found (may be wrong login or password)"))
			if err != nil {
				log1.Error(err.Error())
				break
			}
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
var password = document.getElementById("password");
var login = document.getElementById("login");
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
var send=login.value+"~"+password.value;
print("SEND: " + input.value);
ws.send(input.value+"~"+send);
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
<p><input id="login" type="text" value="login">
<p><input id="password" type="text" value="password">
<p>
<button class="btn btn-secondary mx-auto" id="send">Отправить</button>
</form>
<div class="text-center m-3 pt-3" id="output" style="max-height: 70vh; overflow-y: scroll;"></div>
</body>
</html>
`))
