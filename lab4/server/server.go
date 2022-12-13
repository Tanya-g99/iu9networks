package main

import (
	"bytes"
	"log"
	"os/exec"
	"strings"

	"github.com/gliderlabs/ssh"
	log2 "github.com/mgutz/logxi/v1"
	"golang.org/x/term"
)

func handle(data []string) string {

	name := data[0]
	args := data[1:]

	log2.Info("New command", "name", name, "args", args)

	var buf bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &buf

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

func main() {

	ssh.Handle(func(s ssh.Session) {
		log2.Info("Connection from: ", "user", s.User(), "addr", s.RemoteAddr())

		terminal := term.NewTerminal(s, "> ")
		for {
			line, err := terminal.ReadLine()
			if err != nil {
				break
			}
			answer := handle(strings.Split(line, " "))

			if answer != "" {
				terminal.Write(append([]byte(answer), '\n'))
			}
		}

		log2.Info("terminal closed")
	})

	log2.Info("server start work")

	log.Fatal(ssh.ListenAndServe(":2222", nil))

}
