package postgres

import (
	"GoNews/pkg/storage"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestNew(t *testing.T) {
	err := godotenv.Load("./config_testDB.env")
	if err != nil {
		t.Fatal("Файл конфигурации .env не найден")
	}

	connstr, ok := os.LookupEnv("POSTGRES_CONN_STR_TESTDB")
	if !ok {
		t.Fatal("переменная окружения со строкой подключения к БД не найдена")
	}
	_, err = New(connstr)
	if err != nil {
		t.Fatal(err)
	}
}

// Необходимо очищать тестовую БД перед запуском, как это пофиксить я не знаю((
func TestDB_GoNews(t *testing.T) {
	err := godotenv.Load("../../../config.env")
	if err != nil {
		t.Fatal("Файл конфигурации .env не найден")
	}

	connstr, ok := os.LookupEnv("POSTGRES_CONN_STR_TESTDB")
	if !ok {
		t.Fatal("переменная окружения со строкой подключения к БД не найдена")
	}
	rand.New(rand.NewSource(time.Now().UnixNano()))
	want := []storage.Post{
		{
			Link: strconv.Itoa(rand.Intn(100000)),
		},
		{
			Link: strconv.Itoa(rand.Intn(1000000)),
		},
	}
	newDB, err := New(connstr)
	if err != nil {
		t.Fatal(err)
	}
	err = newDB.AddPosts(want)
	if err != nil {
		t.Error("Публикации не были добавлены", err)
	}
	got, err := newDB.GetPosts(2)
	if got[0].Link != want[0].Link && got[1].Link != want[1].Link {
		t.Errorf("Получили %v и %v, ожидали %v и %v", got[0].Link, got[1].Link, want[0].Link, want[1].Link)
	}
}
