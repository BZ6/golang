package web

import (
  "fmt"
)

func Run() {
  fmt.Println("Сервер запущен по адресу: http://195.133.24.164:8080/")
  HandleRequests()
}
