package hn

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type HNStory struct {
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	ID          int    `json:"id"`
	Kids        []int  `json:"kids"`
	Score       int    `json:"score"`
	Time        int    `json:"time"`
	Title       string `json:"title"`
	Type        string `json:"type"`
	URL         string `json:"url"`
	Host        string
	Rank        int
}

type HNClient struct{}

const apiBaseUrl = "https://hacker-news.firebaseio.com/v0"

func (hnItem HNStory) String() string {
	return fmt.Sprintf("ID: %d, Title: %s, Url: %s, Host: %s\n", hnItem.ID, hnItem.Title, hnItem.URL, hnItem.Host)
}

func NewClient() *HNClient {
	return &HNClient{}
}

func (client *HNClient) GetTopStories() ([]int, error) {
	topStoriesUrl := fmt.Sprintf("%s/topstories.json", apiBaseUrl)
	res, err := http.Get(topStoriesUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var topStoriesIds []int

	jsn, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsn, &topStoriesIds)
	if err != nil {
		return nil, err
	}

	return topStoriesIds, nil
}

func (client *HNClient) GetItemWithID(id int) (HNStory, error) {
	getItemUrl := fmt.Sprintf("%s/item/%d.json", apiBaseUrl, id)

	res, err := http.Get(getItemUrl)
	if err != nil {
		return HNStory{}, err
	}
	defer res.Body.Close()

	var story HNStory

	jsn, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return HNStory{}, err
	}

	err = json.Unmarshal(jsn, &story)
	if err != nil {
		return HNStory{}, err
	}
	domain, err := getDomain(story.URL)
	if err == nil {
		story.Host = domain
	}

	return story, nil
}

func getDomain(urlStr string) (string, error) {
	if urlStr == "" {
		return "", nil
	}
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	parts := strings.Split(u.Hostname(), ".")
	domain := parts[len(parts)-2] + "." + parts[len(parts)-1]
	return domain, nil
}
