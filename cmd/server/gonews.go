package main

import (
	"GoNews/pkg/api"
	"GoNews/pkg/rss"
	"GoNews/pkg/storage"
	"GoNews/pkg/storage/postgres"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type config struct {
	URLS          []string `json:"rss"`
	RequestPeriod int      `json:"request_period"`
}

func main() {
	err := godotenv.Load("./config.env")
	if err != nil {
		log.Fatal("Файл конфигурации .env не найден.")
	}
	connstr, ok := os.LookupEnv("POSTGRES_CONN_STR")
	if !ok {
		log.Fatal("Переменная окружения для подключения к БД не найдена.")
	}
	//Создаем новый экземпляр БД
	db, err := postgres.New(connstr)
	if err != nil {
		log.Fatal(err)
	}
	//Инициализируем API
	api := api.New(db)
	//Читаем файл конфигурации config.json
	b, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatalf("Ошибка чтения файла .json: %v", err)
	}
	var cfg config
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		log.Fatalf("Ошибка расшифровки json: %v", err)
	}
	postsCH := make(chan []storage.Post)
	errCH := make(chan error)

	for _, url := range cfg.URLS {
		go parse(url, db, cfg.RequestPeriod, postsCH, errCH)
	}
	//горутина для записи ошибок
	go func() {
		for err := range errCH {
			log.Printf("Новая ошибка: %f", err)
		}
	}()
	//горутина для записи новостей в БД
	go func() {
		for post := range postsCH {
			_ = db.AddPosts(post)
		}
	}()

	err = http.ListenAndServe(":8080", api.Router())
	if err != nil {
		log.Fatal(err)
	}
}

// Функция для чтения новостей из RSS потока.
func parse(url string, db *postgres.Storage, period int, posts chan<- []storage.Post, errs chan<- error) {

	for {
		news, err := rss.Parse(url)
		if err != nil {
			errs <- err
			continue
		}
		posts <- news
		time.Sleep(time.Duration(period) * time.Minute)
	}
}
