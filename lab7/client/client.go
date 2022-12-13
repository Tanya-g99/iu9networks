package main

import (
	"bufio"
	"fmt"
	"net/smtp"
	"os"

	log "github.com/mgutz/logxi/v1"
)

func main() {
	auth := smtp.PlainAuth("", "t.gnatenko.2003@mail.ru", "df0sc9qYq04kWfg7wkyD",
		"smtp.mail.ru")
	var to string
	fmt.Println("Message To:")
	fmt.Scanln(&to)
	c, err := smtp.Dial("t.gnatenko.2003@mail.ru")
	if err != nil {
		log.Error(err.Error())
	}
	if c.Verify(to) != nil {
		log.Info("нет почты")
	}
	fmt.Println("Message Subject:")
	in := bufio.NewReader(os.Stdin)
	subj, _ := in.ReadString('\n')
	fmt.Println("Message Body:")
	body, _ := in.ReadString('\n')
	fmt.Println("Message Name sender:")
	name, _ := in.ReadString('\n')
	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subj + "\r\n" +
		"\r\n" +
		"Здравствуйте, " + name + body + "\r\n")
	err = smtp.SendMail("smtp.mail.ru:25", auth, "t.gnatenko.2003@mail.ru", []string{to}, msg)
	if err != nil {
		log.Error(err.Error())
	} else {
		log.Info("Письмо успешно отправлено!")
	}
}
