package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/mgutz/logxi/v1"
	"github.com/mmcdole/gofeed"
)

func main() {
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
	for _, item := range feed.Items {
		var isInTable bool
		db.QueryRow("SELECT EXISTS (select * from `Typic` where title = ?)", item.Title).Scan(&isInTable)
		if isInTable {
			log.Info("in table")
			continue
		}

		_, err := db.Exec("insert into `Typic` (`title`, `category`, `description`, `date`, `time`) values (?, ?, ?, ?, ?);",
			item.Title, item.Categories[0], item.Description, item.PublishedParsed.UTC(), item.PublishedParsed.Local())

		if err != nil {
			log.Error(err.Error())
		} else {
			log.Info("add ", item.Title)
		}

	}

}
