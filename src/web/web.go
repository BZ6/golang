package web

import (
  "fmt"
)

func Run() {
  fmt.Println("Сервер запущен по адресу: http://37.77.104.223:8080/")
  HandleRequests()
}
