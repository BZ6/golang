package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/BZ6/golang/src/game"
)

func GamePage(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	var uni game.Universe
	uni.Init(game.WIDTH, game.HEIGHT)

	for {
		output := uni.ShowHTML()
		uni.NextStep()

		if err := conn.WriteMessage(websocket.TextMessage, []byte(output)); err != nil {
			fmt.Println(err)
			return
		}

		time.Sleep(time.Second / game.FPS)
	}
}