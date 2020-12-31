package rest

import "net/http"

type endpoint struct {
	path string
	handler func(writer http.ResponseWriter, request *http.Request)
}

var handlers = []endpoint {
	{
		"/login",
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Login"))
		},
	},
}
