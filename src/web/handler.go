package web

import (
  "net/http"
)

func HandleRequests() {
  http.HandleFunc("/", HomePage)
  http.HandleFunc("/game", GamePage)
  http.ListenAndServe(":8080", nil)
}
