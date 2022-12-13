package main

import (
	"bytes"
	"crypto/tls"
	"database/sql"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/smtp"
	"strings"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/mgutz/logxi/v1"
)

type lett struct {
	Name    string
	Message string
}

type Mail struct {
	Sender  string
	To      []string
	Subject string
	Body    string
}

func BuildMail(mail Mail) []byte {

	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("From: %s\r\n", mail.Sender))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";")))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", mail.Subject))

	fmt.Println("Do you want to attach a file (if so, enter a name, else enter 'no')?")
	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Error(err.Error())
	}

	for _, f := range files {
		fmt.Println(f.Name())
	}

	var NamesFile []string
	var fileName string
	fmt.Scanln(&fileName)
	for fileName != "no" {
		NamesFile = append(NamesFile, fileName)
		fmt.Println("More?")
		fmt.Scanln(&fileName)

	}

	boundary := "my-boundary-779"
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n",
		boundary))

	buf.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
	buf.WriteString("Content-Type: text/html; charset=\"utf-8\"\r\n")
	buf.WriteString(fmt.Sprintf("\r\n%s\r\n", mail.Body))
	if len(NamesFile) != 0 {
		for _, fileName := range NamesFile {

			data, err := ioutil.ReadFile(fileName)

			if err != nil {
				log.Error(err.Error())
			} else {
				buf.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
				buf.WriteString("Content-Type: text/plain; charset=\"utf-8\"\r\n")
				buf.WriteString("Content-Transfer-Encoding: base64\r\n")
				buf.WriteString("Content-Disposition: attachment; filename=" + fileName + "\r\n")
				buf.WriteString("Content-ID: <" + fileName + ">\r\n\r\n")

				b := make([]byte, base64.StdEncoding.EncodedLen(len(data)))
				base64.StdEncoding.Encode(b, data)
				buf.Write(b)
				log.Info("File is attached success", "file: ", fileName)
			}
		}

		buf.WriteString(fmt.Sprintf("\r\n--%s", boundary))
		buf.WriteString("--")
	}

	return buf.Bytes()
}

func main() {
	auth := smtp.PlainAuth("", "t.gnatenko.2003@mail.ru", "RaDhZr4jNuf77KMkm2f6",
		"smtp.mail.ru")
	// var to, message, name, subj string
	// fmt.Println("Message To:")
	// fmt.Scanln(&to)
	// fmt.Println("Message Subject:")
	// fmt.Scanln(&subj)
	// fmt.Println("Message Body:")
	// fmt.Scanln(&message)
	// fmt.Println("Print Name:")
	// fmt.Scanln(&name)

	to := "tanya-g99@ya.ru" //"danila@posevin.com"
	message := "Тупик летел и упал. Логично, ведь это тупик"
	name := "Кто-то"
	subj := "Tupic"

	code := ""
	response := ""

	tlsConfig := tls.Config{
		ServerName:         "smtp.mail.ru",
		InsecureSkipVerify: true,
	}
	log.Info("Establish TLS connection")
	conn, connErr := tls.Dial("tcp", fmt.Sprintf("%s:%d", "smtp.mail.ru", 465),
		&tlsConfig)
	if connErr != nil {
		log.Error(connErr.Error())
	}
	defer conn.Close()
	log.Info("Create new email client")
	client, clientErr := smtp.NewClient(conn, "smtp.mail.ru")
	if clientErr != nil {
		log.Error(clientErr.Error())
	}
	defer client.Close()
	log.Info("Setup authenticate credential")

	if err := client.Auth(auth); err != nil {
		log.Error(err.Error())
	}

	if err := client.Mail(to); err != nil {
		log.Error(err.Error())
		code = (err.Error()[:3])
		response = err.Error()[3 : len(err.Error())-1]
		for code == "550" {
			fmt.Println("Enter new Message To:")
			fmt.Scan(&to)
			if err := client.Rcpt(to); err != nil {
				log.Error(err.Error())
				code = (err.Error()[:3])
				response = err.Error()[3 : len(err.Error())-1]
			}
		}
	}

	body := template.Must(template.New("data").Parse(`
<table bgcolor="#F0FFF0" border="0" cellpadding="0" cellspacing="0"
style="margin:0; padding: 0 10px">
<tr>
<td>
<center style="max-width: 600px; width: 100%;">
<p><b>Здравствуйте, {{.Name}}!</b></p>
<p><i>{{.Message}}</i></p>
</center>
</td>
</tr>
</table>`))
	buf := new(bytes.Buffer)
	body.Execute(buf, lett{
		Name:    name,
		Message: message,
	})
	request := Mail{
		Sender:  "t.gnatenko.2003@mail.ru",
		To:      []string{to},
		Subject: subj,
		Body:    buf.String(),
	}
	msg := BuildMail(request)
	err := smtp.SendMail("smtp.mail.ru:25", auth, "t.gnatenko.2003@mail.ru", []string{to}, msg)
	if err != nil {
		log.Error(err.Error())

	} else {
		log.Info("sending a mail success!")
	}
	db, err := sql.Open("mysql", "iu9networkslabs:Je2dTYr6@tcp(students.yss.su)/iu9networkslabs")

	if err != nil {
		log.Error(err.Error())
		return
	}
	defer db.Close()
	_, err = db.Exec("insert into `lab7` (`Адрес электронной почты получателя`, `Тема сообщения`, `Текст сообщения`, `Имя получателя письма`, `код ответа SMTP сервера`, `расшифровка ответа SMTP сервера`) values (?, ?, ?, ?, ?, ?);",
		// 	to, subj, message, name, "", "")
		to, subj, message, name, code, response)

	if err != nil {
		log.Error(err.Error())
	} else {
		log.Info("add")
	}
}
