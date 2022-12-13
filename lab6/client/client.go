package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/mgutz/logxi/v1"

	"github.com/jlaffaye/ftp"
	"github.com/mmcdole/gofeed"
)

func minHour(i int) string {
	if i > 9 {
		return ""
	} else {
		return "0"
	}
}

func timeStr(currentTime time.Time, str string) string {

	return fmt.Sprintf("%s%d-%d-%d%s%s%d:%s%d:%s%d",
		str,
		currentTime.Day(),
		currentTime.Month(),
		currentTime.Year(),
		str,
		minHour(currentTime.Local().Hour()),
		currentTime.Local().Hour(),
		minHour(currentTime.Local().Minute()),
		currentTime.Local().Minute(),
		minHour(currentTime.Local().Second()),
		currentTime.Local().Second())
}

func main() {
	username := "ftpiu8"
	password := "3Ru7yOTA"
	hostname := "students.yss.su"
	port := "ftp"
	// username := "user"
	// password := "123456"
	// hostname := "localhost"
	// port := "2121"
	//Init client
	c, err := ftp.Dial(hostname + ":" + port)
	if err != nil {
		log.Error("Can't connect to server", "err", err.Error())
		return
	}
	log.Info("Connect to server success")

	err = c.Login(username, password)
	if err != nil {
		log.Error("User or password wrong")
		return
	}
	log.Info("Client entry success")

	// send the commands

	flag := false
	for {
		if flag {
			break
		}
		fmt.Print("> ")
		reader := bufio.NewReader(os.Stdin)
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		cmd := strings.Split(string(line), " ")

		switch cmd[0] {
		case "exit":
			{
				flag = true
			}
		case "help":
			{
				for _, v := range []string{"exit", "post", "get", "make", "delete", "news", "tree\n"} {
					fmt.Print(" ", v)
				}
			}
		case "post":
			{
				file, err := os.Open(cmd[1])
				if err != nil {
					log.Error(err.Error())
					os.Exit(1)
				}
				defer file.Close()

				err = c.Stor(cmd[1], file)
				if err != nil {
					panic(err)
				}
				log.Info("Posting success")
			}
		case "get":
			{
				r, err := c.Retr(cmd[1])
				if err != nil {
					panic(err)
				}
				defer r.Close()
				buf, err := ioutil.ReadAll(r)
				if err != nil {
					log.Error(err.Error())
					os.Exit(1)
				}
				file, err := os.Create(cmd[1])

				if err != nil {
					log.Error(err.Error())
					os.Exit(1)
				}
				defer file.Close()
				file.WriteString(string(buf))
				log.Info("Getting success")
			}
		case "make":
			{
				err := c.MakeDir(cmd[1])
				if err != nil {
					log.Error(err.Error())
				}
			}
		case "delete":
			{
				err := c.Delete(cmd[1])
				if err != nil {
					log.Error(err.Error())
				}
				log.Info("Delete success")
			}
		case "tree":
			{
				dir, err := c.CurrentDir()
				if err != nil {
					log.Error(err.Error())
				} else {
					answer, err := c.NameList(dir)
					if err != nil {
						log.Error(err.Error())
					} else {
						for _, v := range answer {
							fmt.Println(v)
						}
					}
				}
			}
		case "news":
			{
				db, err := sql.Open("mysql", "iu9networkslabs:Je2dTYr6@tcp(students.yss.su)/iu9networkslabs")

				if err != nil {
					log.Error(err.Error())
					return
				}
				defer db.Close()

				fp := gofeed.NewParser()
				feed, err := fp.ParseURL("https://news.mail.ru/rss/90/")
				if err != nil {
					log.Error(err.Error())
				}

				// make file
				fileName := "Gnatenko_Tanya" + timeStr(time.Now(), "_") + ".txt"
				count := 0
				buf := ""
				for _, item := range feed.Items {
					if count == 5 {
						break
					}
					var isInTable bool
					db.QueryRow("SELECT EXISTS (select * from `Typic` where title = ?)", item.Title).Scan(&isInTable)
					if isInTable {
						continue
					}
					_, err = db.Exec("insert into `Typic` (`title`, `category`, `description`, `date`, `time`) values (?, ?, ?, ?, ?);",
						item.Title, item.Categories[0], item.Description, item.PublishedParsed.UTC(), item.PublishedParsed.Local())

					if err != nil {
						log.Error(err.Error())
					} else {
						count++
						for _, v := range []string{item.Title, item.Categories[0], item.Description, timeStr(item.PublishedParsed.UTC(), " "), timeStr(item.PublishedParsed.Local(), " "), "\n"} {
							buf += " " + v
						}

					}

				}

				if buf != "" {
					r := strings.NewReader(buf)
					c.Append(fileName, r)
					log.Info("Posting news successful")
				} else {
					log.Error("Posting news unsuccessful")
				}
			}
		}
	}

	// Do something with the FTP conn

	if err := c.Quit(); err != nil {
		log.Error(err.Error())
	}
}
