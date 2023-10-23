package api

import "net/http"

func (api *API) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json") //Устанавливаем заголовки ответа сервера
		next.ServeHTTP(res, req)
	})
}
