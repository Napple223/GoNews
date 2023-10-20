package storage

//пакет storage предназначен для отделения реализации логики работы приложения от конкретной БД

//Структура данных поста
type Post struct {
	ID      int    //id публикации
	Title   string //заголовок публикации
	Content string //тело публикации
	PubTime int64  //время публикации
	Link    string //ссылка на источник
}

//Интерфейс - контракт для работы с БД
type Interface interface {
	GetPosts(int) ([]Post, error) //Вывод заданного количества постов из БД
	AddPosts([]Post) error        //Добавление поста в БД
}
