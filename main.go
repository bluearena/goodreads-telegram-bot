package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	goodreads "github.com/halink0803/goodreads-telegram-bot/goodreads"
	tb "gopkg.in/tucnak/telebot.v2"
)

type BotConfig struct {
	Key string `json:"bot_key"`
}

var currentCommand string

func readConfigFromFile(path string) (BotConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return BotConfig{}, err
	} else {
		result := BotConfig{}
		err := json.Unmarshal(data, &result)
		return result, err
	}
}

func main() {
	path := "./config.json"
	botKey, _ := readConfigFromFile(path)
	bot, err := tb.NewBot(tb.Settings{
		Token:  botKey.Key,
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
		handleSearch(bot, m)
	})

	bot.Handle("/list", func(m *tb.Message) {
		getListOfBookShelves(bot, m)
	})

	bot.Handle(tb.OnText, func(m *tb.Message) {
		handleText(bot, m)
	})

	bot.Start()
}

func getListOfBookShelves(bot *tb.Bot, m *tb.Message) {
	result, _ := goodreads.GetListShelves()
	log.Printf("list result: %+v", result)
	inlineKeys := [][]tb.InlineButton{}
	for _, v := range result.Shelves.UserShelves {
		inlineBtn := tb.InlineButton{
			Unique: fmt.Sprintf("sad_moon_%d", v.ID),
			Text:   v.Name,
		}
		inlineKeys = append(inlineKeys, []tb.InlineButton{inlineBtn})
	}
	bot.Send(m.Sender, "Message", &tb.ReplyMarkup{
		InlineKeyboard: inlineKeys,
	})
}

func handleSearch(bot *tb.Bot, m *tb.Message) {
	queries := strings.Split(m.Text, " ")
	if len(queries) > 1 {
		query := strings.Join(queries[1:], " ")
		searchBook(bot, query, m)
	} else {
		currentCommand = "search"
		bot.Send(m.Sender, "Which book do you want to search")
	}
}

func handleText(bot *tb.Bot, m *tb.Message) {
	switch currentCommand {
	case "search":
		searchBook(bot, m.Text, m)
		currentCommand = ""
	}
}

func searchBook(bot *tb.Bot, query string, m *tb.Message) {
	result, _ := goodreads.SearchBook(query)
	for i := 0; i < 5 && i < len(result.Books); i++ {
		book := result.Books[i]
		var message string
		message += fmt.Sprintf("Title: %s\n", book.Title)
		message += fmt.Sprintf("Author: %s\n", book.Author)
		message += fmt.Sprintf("Average rating: %.2f\n", book.AverageRating)
		message += fmt.Sprintf("https://goodreads.com/book/show/%d", book.ID)
		bot.Send(m.Sender, message)
	}
}
