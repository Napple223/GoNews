package rss

import (
	"GoNews/pkg/storage"
	"encoding/xml"
	"io"
	"net/http"
	"strings"
	"time"

	strip "github.com/grokify/html-strip-tags-go"
)

// Структура RSS потока
type Stream struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	Version string   `xml:"version,attr"`
	Dc      string   `xml:"dc,attr"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Text           string `xml:",chardata"`
	Title          string `xml:"title"`
	Description    string `xml:"description"`
	Language       string `xml:"language"`
	ManagingEditor string `xml:"managingEditor"`
	Generator      string `xml:"generator"`
	Image          struct {
		Text  string `xml:",chardata"`
		Link  string `xml:"link"`
		URL   string `xml:"url"`
		Title string `xml:"title"`
	} `xml:"image"`
	Link  string `xml:"link"`
	Items []Item `xml:"item"`
}

type Item struct {
	Text     string   `xml:",chardata"`
	Title    string   `xml:"title"`
	Link     string   `xml:"link"`
	Content  string   `xml:"description"`
	PubTime  string   `xml:"pubDate"`
	Guid     string   `xml:"guid"`
	Creator  string   `xml:"creator"`
	Category []string `xml:"category"`
}

// Считывает RSS поток и возвращает слайс раскодированных новостей
func Parse(url string) ([]storage.Post, error) {
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
		p.Content = strip.StripTags(items.Content)
		nTime := strings.ReplaceAll(items.PubTime, ",", "")
		t, err := time.Parse("Mon 2 Jan 2006 15:04:05 -0700", nTime)
		if err != nil {
			t, err = time.Parse("Mon 2 Jan 2006 15:04:05 GMT", nTime)
		}
		if err == nil {
			p.PubTime = t.Unix()
		}
		p.Link = items.Link
		posts = append(posts, p)
	}
	return posts, nil
}
