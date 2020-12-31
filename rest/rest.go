package rest

import (
	"log"
	"net/http"
	"strconv"
)

const Port = 8080

func StartService() {
	log.Printf("%d endpoints were registered", registerEndpoints(handlers))
	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(Port), nil))
}

func registerEndpoints(points []endpoint) int {
	var amount int
	var p endpoint
	for amount, p = range points {
		http.HandleFunc(p.path, p.handler)
		log.Printf("Registered #%d for path '%v'", amount + 1, p.path)
	}
	return amount + 1
}