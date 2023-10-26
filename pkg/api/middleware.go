package api

import "net/http"

func (api *API) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		//корректно работает только в мозиле, если убрать Content-Type, то ничего не меняется
		res.Header().Set("Content-Type", "text/html") //Устанавливаем заголовки ответа сервера
		res.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(res, req)
	})
}
