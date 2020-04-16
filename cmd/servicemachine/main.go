package main

import (
	//"encoding/json"
	//"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/{type}", typeHandle).Methods("GET")
	http.ListenAndServe(":8989", router)
}
