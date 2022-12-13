---
предмет: Компьютерные сети
название: Разработка SSH-сервера и SSH-клиента
номер: 4
тип_работы: Лабораторная работа
группа: ИУ9-32Б
автор: Гнатенко Т. А.
преподаватель: Посевин Д. П.
---

# Цели

Изучить принципы работы протокола ssh

# Часть 1. Разработка SSH-сервера

## Задачи

Реализовать ssh сервер на языке GO с применением указанных пакетов и
запустить его на localhost. Проверка работы должна проводиться путем использования
программы ssh в ОС Linux/Unix или PuTTY в ОС Windows. Должны работать следующие
функции:
◦ авторизация клиента на ssh сервере;
◦ создание директории на удаленном сервере;
◦ удаление директории на удаленном сервере;
◦ вывод содержимого директории;
◦ перемещение файлов из одной директории в другую;
◦ удаление файла по имени;
◦ вызов внешних приложений, например ping.

## Решение

### Исходный код

**`server.go`**

```go

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

```

## Вывод

![](pic/Screenshot%20from%202022-10-04%2016-37-03.png)

![](pic/Screenshot%20from%202022-10-04%2016-35-25.png)

# Часть 2. Разработка SSH-клиента

## Задачи

Реализовать ssh-клиент и запустить его на localhost.
Протестировать соединение Go SSH-клиента к серверу реализованному в
предыдущей задаче, а также к произвольному ssh серверу.
Требования: SSH-клиент должен поддерживать следующие функции:
◦ авторизация клиента на SSH-сервере;
◦ создание директории на удаленном SSH-сервере;
◦ удаление директории на удаленном SSH-сервере;
◦ вывод содержимого директории;
◦ перемещение файлов из одной директории в другую;
◦ удаление файла по имени;
◦ вызов внешних приложений, например ping.

## Решение

### Исходный код

```go

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
	// username := "user"
	// password := "S4EmRoIhDy/w1bIzlOU51lphBN4=|+5jIYH84ai5UcFexe6VSFCaMQV0= ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC10H6D+szTfU4FZvGlzUNYRYfrSNGPfOSHP92zcmiex7y5rQmx3FQoVDer3ClwgsesAW9VUKMI2Nmweo7NRUXV6uZ/0M5lr7VQtJ1MgopcXQQdY5S5MQiIo10rWFN5YyRwNIw48g7/AZZXDheGiyykhMM+BODgpB7ivqlPZlcOmMGgu3ULIUbaAxTDeIsE0jbtAKkoYEMDFGRS0txFM0uj2T5HwVV7jLqcjxCjSf7E5UPRAqQeOoztqRGszKAzdGVV4lWbDpNxg2cuIXcrEX2lFs9wQlDGWdrofj6J/zmIx3kcvBGrxGwO0lZQLDu2EQguam6iXI9wjQDHOW3Y5Lml"
	// hostname := "localhost"
	// port := "2222"

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

```

## Вывод

![](pic/Screenshot%20from%202022-10-04%2019-37-53.png)


