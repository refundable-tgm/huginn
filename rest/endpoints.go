package rest

import "net/http"

type endpoint struct {
	path string
	handler func(writer http.ResponseWriter, r *http.Request)
}

var handlers = []endpoint {
	{
		//tba
	},
	{
		//tba
	},
}
