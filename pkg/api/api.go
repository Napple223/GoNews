package api

import (
	"GoNews/pkg/storage"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	errSIS = http.StatusInternalServerError
)

// Программный интерфейс сервера
type API struct {
	db     storage.Interface
	router *mux.Router
}

// Регистрация API в маршрутизаторе запросов
func (api *API) endpoints() {
	api.router.Use(api.Middleware)
	api.router.HandleFunc("/posts/{n}", api.postsHandler).Methods(http.MethodGet, http.MethodOptions)
	api.router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
}

// Конструктор объекта API
func New(db storage.Interface) *API {
	api := API{
		db: db,
	}
	api.router = mux.NewRouter()
	api.endpoints()
	return &api
}

// Получение маршрутизатора запросов
func (api *API) Router() *mux.Router {
	return api.router
}

// Получение заданного количества новостей
func (api *API) postsHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	if req.Method == http.MethodOptions {
		return
	}
	str := mux.Vars(req)["n"]
	n, err := strconv.Atoi(str)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	posts, err := api.db.GetPosts(n)
	if err != nil {
		http.Error(res, err.Error(), errSIS)
	}
	json.NewEncoder(res).Encode(posts)
}
