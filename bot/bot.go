package bot

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)


type BotConfig struct {
	Key string `json:"bot_key"`
}

func GetConfigFromFile(path string) (BotConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return BotConfig{}, err
	} else {
		result := BotConfig{}
		err := json.Unmarshal(data, &result)
		return result, err
	}
}

func NewBot(key string) *Bot {
}

func (self *Bot) GetResponse(
	method string, req_url string,
	params map[string]string) ([]byte, error) {

	client := &http.Client{
		Timeout: time.Duration(30 * time.Second),
	}
	req, _ := http.NewRequest(method, req_url, nil)
	req.Header.Add("Accept", "application/json")

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	var err error
	var resp_body []byte
	resp, err := client.Do(req)
	if err != nil {
		return resp_body, err
	} else {
		defer resp.Body.Close()
		resp_body, err = ioutil.ReadAll(resp.Body)
		log.Printf("request to %s, got response: %s\n", req.URL, resp_body)
		return resp_body, err
	}
}

func (self *Bot) SetWebHook(url string) {
	resp_body, err := self.GetResponse(
		"GET",
		TELEGRAM_BOT_API+self.Key+"/setWebhook",
		map[string]string{
			"url": url,
		},
	)
	if err != nil 
}
