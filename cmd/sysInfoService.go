package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/user"

	"github.com/gorilla/mux"
	"github.com/zcalusic/sysinfo"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/{category}", categoryHandle).Methods("GET")
	http.ListenAndServe(":9999", r)
}

func categoryHandle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	current, err := user.Current()
	if err != nil {
		log.Fatal("Error getting current user")
	}

	if current.Uid != "0" {
		log.Fatal("Requires superuser privilege")
	}

	var si sysinfo.SysInfo
	si.GetSysInfo()
	var data []byte
	switch vars["category"] {
	case "sysinfo":
		data, err = json.MarshalIndent(&si.Meta, "", "  ")
	case "node":
		data, err = json.MarshalIndent(&si.Node, "", "  ")
	case "os":
		data, err = json.MarshalIndent(&si.OS, "", "  ")
	case "kernel":
		data, err = json.MarshalIndent(&si.Kernel, "", "  ")
	case "product":
		data, err = json.MarshalIndent(&si.Product, "", "  ")
	case "board":
		data, err = json.MarshalIndent(&si.Board, "", "  ")
	case "chassis":
		data, err = json.MarshalIndent(&si.Chassis, "", "  ")
	case "bios":
		data, err = json.MarshalIndent(&si.BIOS, "", "  ")
	case "cpu":
		data, err = json.MarshalIndent(&si.CPU, "", "  ")
	case "memory":
		data, err = json.MarshalIndent(&si.Memory, "", "  ")
	case "storage":
		data, err = json.MarshalIndent(&si.Storage, "", "  ")
	case "network":
		data, err = json.MarshalIndent(&si.Network, "", "  ")
	default:
		data, err = json.MarshalIndent(&si, "", "  ")
	}
	if err != nil {
		log.Fatal("error to display info")
	}
	fmt.Fprintf(w, string(data))
	w.WriteHeader(http.StatusOK)
}
