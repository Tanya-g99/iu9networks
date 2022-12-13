package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

const index_html = `
<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<link
	href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.2/dist/css/bootstrap.min.css"
    rel="stylesheet"
    integrity="sha384-EVSTQN3/azprG1Anm3QDgpJLIm9Nao0Yz1ztcQTwFspd3yD65VohhpuuCOmLASjC"
	crossorigin="anonymous"
/>
<title>Form</title>
</head>
<body>
<form class="text-center m-3" method="POST" action="{{.}}">
<div>
<input name="variables" type="text">
</div>
<button class="btn btn-secondary m-3" type="submit" value="submit">Отправить</button>
</form>
</body>
</html>
`

var addr = flag.String("addr", "localhost:8082", "http service address")

func helloHandler(w http.ResponseWriter, r *http.Request) {
	//str := strings.Split(r.URL.Path, "/")
	switch r.Method {
	case "GET":
		tmpl := template.Must(template.New("data").Parse(index_html))
		tmpl.Execute(w, "./")
	case "POST":
		if err := r.ParseForm(); err != nil {
			fmt.Fprintf(w, "ParseForm() err: %v", err)
			return
		}
		flag.Parse()
		log.SetFlags(0)
		u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
		log.Printf("connecting to %s", u.String())
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Fatal("dial: ", err)
		}
		defer c.Close()
		x := r.FormValue("variables")
		err = c.WriteMessage(websocket.TextMessage, []byte(x))
		if err != nil {
			log.Println("write: ", err)
			return
		}
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("Answer: %s", message)
		fmt.Fprint(w, string(message))
	}
}
func main() {
	log.Println("Starting http server at port 8080")
	http.HandleFunc("/", helloHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
