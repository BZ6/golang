package web

import (
    "net/http"
)

func HandleRequests() {
    http.HandleFunc("/", HomePage)
    http.ListenAndServe(":8080", nil)
}