package rest

import (
	"log"
	"net/http"
	"strconv"
)

const Port = 42069

func StartService() {
	registerEndpoints(handlers)
	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(Port), nil))
}

func registerEndpoints(points []endpoint)  {
	for i, p := range points {
		http.HandleFunc(p.path, p.handler)
		log.Printf("Registered #%d for path '%v'", i, p.path)
	}
}