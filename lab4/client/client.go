package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

func main() {

	username := "test"
	password := "SDHBCXdsedfs222"
	hostname := "151.248.113.144"
	port := "443"

	// SSH client config
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		// Non-production only
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to host
	client, err := ssh.Dial("tcp", hostname+":"+port, config)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// Create sesssion
	sesssion, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer sesssion.Close()

	// StdinPipe for commands
	stdin, err := sesssion.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	// Enable system stdout
	// Comment these if you uncomment to store in variable
	sesssion.Stdout = os.Stdout
	sesssion.Stderr = os.Stderr

	// Start remote shell
	err = sesssion.Shell()
	if err != nil {
		log.Fatal(err)
	}

	// send the commands
	for {
		reader := bufio.NewReader(os.Stdin)
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		cmd := string(line)
		if cmd != "exit" {
			_, err = fmt.Fprintf(stdin, "%s\n", cmd)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			break
		}

	}

	// Wait for sesssion to finish
	err = sesssion.Wait()
	if err != nil {
		log.Fatal(err)
	}

}
