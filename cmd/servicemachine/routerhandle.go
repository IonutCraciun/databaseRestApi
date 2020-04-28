package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func typeHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var data []byte
	var err error
	switch vars["type"] {
	case "cpu":
		//data, err = json.MarshalIndent(GetCpuDetails(), "", "    ") Usage without goroutine
		routine := GetCpuDetailsRoutine()
		data, err = json.MarshalIndent(<-routine, "", "    ")
	case "storage":

		data, err = json.MarshalIndent(GetStorageDetails(), "", "    ")
	case "memory":
		// different approach, same behaviour
		var blockMemory BlockMemory
		blockMemory.GetMemoryDetails()
		data, err = json.MarshalIndent(blockMemory, "", "    ")
	}
	if err != nil {
		log.Fatal("error to display info")
	}
	fmt.Println(string(data))
	fmt.Fprintf(w, string(data))
	w.WriteHeader(http.StatusOK)

}
