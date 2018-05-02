package main

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/oauth2"
)

type GoodReads struct {
	Key       string
	SecretKey string
	UserID    string
}

func NewGoodReads(key, secretKey, userID string) *GoodReads {
	return &GoodReads{
		key,
		secretKey,
		userID,
	}
}

type GoodReadsConfig struct {
	Key       string `json:"goodreads_key"`
	SecretKey string `json:"goodreads_secet_key"`
	UserID    string `json:"user_id"`
}

func GetConfigFromFile(path string) (GoodReadsConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return GoodReadsConfig{}, err
	} else {
		result := GoodReadsConfig{}
		err := json.Unmarshal(data, &result)
		return result, err
	}
}

func (self GoodReads) GetResponse(
	method string, req_url string,
	params map[string]string) ([]byte, error) {

	// config oauth
	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     self.Key,
		ClientSecret: self.SecretKey,
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://www.goodreads.com/oauth/request_token",
			AuthURL:  "https://www.goodreads.com/oauth/access_token",
		},
	}

	log.Printf("Oauth: %+v", conf)

	httpClient := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}
	req, err := http.NewRequest(method, req_url, nil)
	if err != nil {
		log.Println(err.Error())
	}
	req.Header.Add("Accept", "application/json")

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	var resp_body []byte

	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog: %v", url)

	ctx = context.WithValue(ctx, oauth2.HTTPClient, httpClient)
	tok, err := conf.Exchange(ctx, "authorization-code")
	if err != nil {
		log.Println(err.Error())
	}
	client := conf.Client(ctx, tok)
	log.Printf("Client %+v", client)
	resp, err := client.Do(req)
	if err != nil {
		return resp_body, err
	}
	defer resp.Body.Close()
	resp_body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
	}
	log.Printf("request to %s, got response: %s\n", req.URL, resp_body)
	return resp_body, err

}

func (self GoodReads) GetListShelves() (GoodreadsUserShelves, error) {
	resp_body, err := self.GetResponse(
		"GET",
		"https://www.goodreads.com/shelf/list.xml",
		map[string]string{
			"key":     self.Key,
			"user_id": self.UserID,
		},
	)
	result := GoodreadsUserShelves{}
	if err == nil {
		xml.Unmarshal(resp_body, &result)
	} else {
		log.Println(err.Error())
	}
	return result, err
}

func (self GoodReads) SearchBook(searchQuery string) (BookSearchResponse, error) {
	resp_body, err := self.GetResponse(
		"GET",
		"https://www.goodreads.com/search/index.xml",
		map[string]string{
			"key": self.Key,
			"q":   searchQuery,
		},
	)
	result := BookSearchResponse{}
	if err == nil {
		err := xml.Unmarshal(resp_body, &result)
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		log.Println(err.Error())
	}
	return result, err
}

func (self GoodReads) AddBookToShelf(bookId int, shelfName string) error {
	_, err := self.GetResponse(
		"POST",
		"https://www.goodreads.com/shelf/add_to_shelf.xml",
		map[string]string{
			"name":    shelfName,
			"book_id": strconv.Itoa(bookId),
		},
	)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return err
}
