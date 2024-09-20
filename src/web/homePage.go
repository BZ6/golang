package web

import (
  "net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
  http.ServeFile(w, r, "pages/index.html")
}
