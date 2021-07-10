package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/ramin0x53/telegram_alertbot/src/price"
	tb "gopkg.in/tucnak/telebot.v2"
)

var TKN string = "1839193222:AAGAtkSI2r6Z9WWWwrp4QtKehBWlyd3N404"

type sendMessageReqBody struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

//c = current price    p = price in database     r = rule in database
func comp(c string, p string, r string) bool {
	c1, _ := strconv.ParseFloat(c, 8)
	p1, _ := strconv.ParseFloat(p, 8)

	if r == ">" {
		if c1 >= p1 {
			return true
		} else {
			return false
		}
	} else if r == "<" {
		if c1 <= p1 {
			return true
		} else {
			return false
		}
	} else {
		fmt.Println("Wronge rule !")
		return false
	}
}

func alert() {
	database, _ := sql.Open("sqlite3", "./CryptoTable.db")
	for {
		rows, _ := database.Query("SELECT chatid, crypto, rule, price FROM table1 ")

		var chatid int
		var crypto string
		var rule string
		var price1 string

		for rows.Next() {
			rows.Scan(&chatid, &crypto, &rule, &price1)
			if comp(price.GetPrice(crypto), price1, rule) {
				reqBody := &sendMessageReqBody{
					ChatID: int64(chatid),
					Text:   crypto + " " + rule + " " + price1 + " ⁉️",
				}

				reqBytes, _ := json.Marshal(reqBody)

				http.Post("https://api.telegram.org/bot"+TKN+"/sendMessage", "application/json", bytes.NewBuffer(reqBytes))

				statement, _ := database.Prepare("DELETE FROM table1 WHERE chatid = " + strconv.Itoa(chatid) + " AND crypto = '" + crypto + "' AND rule = '" + rule + "' AND price = '" + price1 + "'")
				_, err := statement.Exec()

				if err != nil {
					log.Fatal(err)
				}

				time.Sleep(300 * time.Millisecond)

			}
		}
	}
}

func main() {
	b, err := tb.NewBot(tb.Settings{

		Token:  TKN,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	database, _ := sql.Open("sqlite3", "./CryptoTable.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS table1 (id INTEGER PRIMARY KEY,chatid INTEGER, crypto TEXT, rule TEXT, price TEXT)")
	statement.Exec()

	go alert()

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
