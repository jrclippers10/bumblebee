package main

import (
    "github.com/gorilla/mux"
    "log"
    "net/http"
)

func StartServer() {
    r := mux.NewRouter()
    r.HandleFunc("/", IndexHandler)
    log.Fatal(http.ListenAndServe(":8000", r))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Bitcoin!\n"))
}