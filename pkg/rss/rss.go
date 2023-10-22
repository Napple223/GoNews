package rss

import (
	"GoNews/pkg/storage"
	"encoding/xml"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
)

// Структура RSS потока
type Stream struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	//Title       string `xml:"title"`
	//Description string `xml:"description"`
	//Link        string `xml:"link"`
	Items []Item `xml:"item"`
}

type Item struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	Content string `xml:"description"`
	PubTime string `xml:"pubData"`
}

// Считывает RSS поток и возвращает слайс раскодированных новостей
func Parse(url string) ([]storage.Post, error) {
	policy := bluemonday.UGCPolicy()
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var s Stream
	err = xml.Unmarshal(b, &s)
	if err != nil {
		return nil, err
	}

	var posts []storage.Post
	for _, items := range s.Channel.Items {
		var p storage.Post
		p.Title = items.Title
		p.Content = policy.Sanitize(items.Content)
		nTime := strings.ReplaceAll(items.PubTime, ",", "")
		t, err := time.Parse("Sun 3 Feb 2013 22:22:0 +0300", nTime)
		if err != nil {
			t, err = time.Parse("Sun 3 Feb 2013 15:20:0 GMT", nTime)
		}
		if err == nil {
			p.PubTime = t.Unix()
		}
		p.Link = items.Link
		posts = append(posts, p)
	}
	return posts, nil
}
