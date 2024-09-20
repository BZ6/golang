package web

import (
  "fmt"
)

func Run() {
  fmt.Println("Сервер запущен по адресу: http://localhost:8080/")
  HandleRequests()
}
