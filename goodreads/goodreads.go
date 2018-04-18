package goodreads

import (
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	GOODREADS_API_KEY string = "tPChp05p5D2dfR0AP6XXyg"
	USER_ID           string = "27284337"
)

type GoodReads struct {
	Key    string
	UserID string
}

func NewGoodReads(key, userID string) *GoodReads {
	return &GoodReads{
		key,
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

func GetResponse(
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

func GetListShelves() (GoodreadsUserShelves, error) {
	resp_body, err := GetResponse(
		"GET",
		"https://www.goodreads.com/shelf/list.xml?key="+GOODREADS_API_KEY+"&user_id="+USER_ID,
		map[string]string{},
	)
	result := GoodreadsUserShelves{}
	if err == nil {
		xml.Unmarshal(resp_body, &result)
	} else {
		log.Println(err.Error())
	}
	return result, err
}

func SearchBook(searchQuery string) (BookSearchResponse, error) {
	resp_body, err := GetResponse(
		"GET",
		"https://www.goodreads.com/search/index.xml",
		map[string]string{
			"key": GOODREADS_API_KEY,
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
