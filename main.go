package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	goodreads "github.com/halink0803/goodreads-telegram-bot/goodreads"
	tb "gopkg.in/tucnak/telebot.v2"
)

type BotConfig struct {
	Key             string `json:"bot_key"`
	GoodreadsKey    string `json:"goodreads_api_key"`
	GoodreadsUserID string `json:"goodreads_user_id"`
}

type Bot struct {
	bot            *tb.Bot
	goodreads      goodreads.GoodReads
	currentCommand string
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
	botConfig, _ := readConfigFromFile(path)
	tbot, err := tb.NewBot(tb.Settings{
		Token:  botConfig.Key,
		Poller: &tb.LongPoller{Timeout: 5 * time.Second},
	})
	mybot := Bot{
		bot: tbot,
		goodreads: goodreads.GoodReads{
			Key:    botConfig.GoodreadsKey,
			UserID: botConfig.GoodreadsUserID,
		},
	}
	if err != nil {
		log.Panic(err)
		return
	}
	mybot.bot.Handle("/start", func(m *tb.Message) {
		log.Println("receive message: ", m.Text)
		mybot.bot.Send(m.Sender, "Thank you for trying")
	})

	mybot.bot.Handle("/search", func(m *tb.Message) {
		mybot.handleSearch(m)
	})

	mybot.bot.Handle("/list", func(m *tb.Message) {
		mybot.getListOfBookShelves(m)
	})

	mybot.bot.Handle(tb.OnText, func(m *tb.Message) {
		mybot.handleText(m)
	})

	mybot.bot.Start()
}

func (self Bot) getListOfBookShelves(m *tb.Message) {
	result, _ := self.goodreads.GetListShelves()
	log.Printf("list result: %+v", result)
	inlineKeys := [][]tb.InlineButton{}
	for _, v := range result.Shelves.UserShelves {
		inlineBtn := tb.InlineButton{
			Unique: fmt.Sprintf("sad_moon_%d", v.ID),
			Text:   v.Name,
		}
		inlineKeys = append(inlineKeys, []tb.InlineButton{inlineBtn})
	}
	self.bot.Send(m.Sender, "Message", &tb.ReplyMarkup{
		InlineKeyboard: inlineKeys,
	})
}

func (self Bot) handleSearch(m *tb.Message) {
	queries := strings.Split(m.Text, " ")
	if len(queries) > 1 {
		query := strings.Join(queries[1:], " ")
		self.searchBook(query, m)
	} else {
		currentCommand = "search"
		self.bot.Send(m.Sender, "Which book do you want to search", &tb.ReplyMarkup{
			ReplyKeyboard: nil,
		})
	}
}

func (self Bot) handleText(m *tb.Message) {
	switch currentCommand {
	case "search":
		self.searchBook(m.Text, m)
		currentCommand = ""
	}
}

func (self Bot) searchBook(query string, m *tb.Message) {
	result, _ := self.goodreads.SearchBook(query)
	for i := 0; i < 5 && i < len(result.Books); i++ {
		book := result.Books[i]
		var message string
		message += fmt.Sprintf("Title: %s\n", book.Title)
		message += fmt.Sprintf("Author: %s\n", book.Author)
		message += fmt.Sprintf("Average rating: %.2f\n", book.AverageRating)
		message += fmt.Sprintf("https://goodreads.com/book/show/%d", book.ID)

		inlineKeys := [][]tb.InlineButton{}
		inlineBtn := tb.InlineButton{
			Unique: fmt.Sprintf("%s_%d", "shelf", book.ID),
			Text:   "Add to shelf",
		}

		self.bot.Handle(&inlineBtn, func(c *tb.Callback) {
			id, _ := strconv.Atoi(strings.Split(inlineBtn.Unique, "_")[1])
			log.Printf("book: %d", id)
			self.handleAddToShelf(id, m)
			self.bot.Respond(c, &tb.CallbackResponse{})
		})

		inlineKeys = append(inlineKeys, []tb.InlineButton{inlineBtn})
		self.bot.Send(m.Sender, message, &tb.ReplyMarkup{
			InlineKeyboard: inlineKeys,
		})

	}
}

func (self Bot) handleAddToShelf(bookID int, m *tb.Message) {
	result, _ := self.goodreads.GetListShelves()
	replyKeys := [][]tb.ReplyButton{}
	for _, v := range result.Shelves.UserShelves {
		replyBtn := tb.ReplyButton{
			Text: v.Name,
		}
		replyKeys = append(replyKeys, []tb.ReplyButton{replyBtn})
	}
	replyBtn := tb.ReplyButton{
		Text: "Add new shelf",
	}
	replyKeys = append(replyKeys, []tb.ReplyButton{replyBtn})
	self.bot.Send(m.Sender, "Which self do you want to set to?", &tb.ReplyMarkup{
		ReplyKeyboard: replyKeys,
	})
}
