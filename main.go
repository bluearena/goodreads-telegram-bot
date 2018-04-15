package main

import (
	"encoding/json"
	"log"
	"time"

	goodreads "github.com/halink0803/goodreads-telegram-bot/goodreads"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	botKey := ""
	bot, err := tb.NewBot(tb.Settings{
		Token:  botKey,
		Poller: &tb.LongPoller{Timeout: 5 * time.Second},
	})
	if err != nil {
		log.Panic(err)
		return
	}
	bot.Handle("/start", func(m *tb.Message) {
		log.Println("receive message: ", m.Text)
		bot.Send(m.Sender, "Thank you for trying")
	})

	bot.Handle("/search", func(m *tb.Message) {
		log.Println("")
	})

	bot.Handle("/list", func(m *tb.Message) {
		getListOfBookShelves()
	})

	bot.Start()
}

func getListOfBookShelves() {
	result, _ := goodreads.GetListShelves()
	data, _ := json.Marshal(result)
	log.Println(string(data))
	// bot.Send(m.Sender, string(data))
}
