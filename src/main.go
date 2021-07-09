package main

import (
	"database/sql"
	"log"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	// "github.com/ramin0x53/telegram_alertbot/src/price"
	"github.com/ramin0x53/telegram_alertbot/src/price"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	b, err := tb.NewBot(tb.Settings{

		Token:  "1839193222:AAGAtkSI2r6Z9WWWwrp4QtKehBWlyd3N404",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	database, _ := sql.Open("sqlite3", "./CryptoTable.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS table1 (id INTEGER PRIMARY KEY,chatid INTEGER, crypto TEXT, rule TEXT, price TEXT)")
	statement.Exec()

	b.Handle("/add", func(m *tb.Message) {
		split := strings.Split(strings.ReplaceAll(m.Text, "/add ", ""), " ")
		statement, _ = database.Prepare("INSERT INTO table1 (chatid, crypto, rule, price) VALUES (?, ?, ?, ?)")
		statement.Exec(strconv.Itoa(m.Sender.ID), split[0], split[1], split[2])
	})

	b.Handle("/rm", func(m *tb.Message) {
		split := strings.Split(strings.ReplaceAll(m.Text, "/rm ", ""), " ")
		statement, _ = database.Prepare("DELETE FROM table1 WHERE chatid = " + strconv.Itoa(m.Sender.ID) + " AND crypto = '" + split[0] + "' AND rule = '" + split[1] + "' AND price = '" + split[2] + "'")
		_, err := statement.Exec()

		if err != nil {
			log.Fatal(err)
		}
	})

	b.Handle("/show", func(m *tb.Message) {
		rows, _ := database.Query("SELECT crypto, rule, price FROM table1 WHERE chatid = " + strconv.Itoa(m.Sender.ID))

		var crypto string
		var rule string
		var price string

		for rows.Next() {
			rows.Scan(&crypto, &rule, &price)
			b.Send(m.Sender, crypto+" "+rule+" "+price)
		}
	})

	b.Handle(tb.OnText, func(m *tb.Message) { b.Send(m.Chat, price.GetPrice(m.Text)) })
	b.Start()
}
