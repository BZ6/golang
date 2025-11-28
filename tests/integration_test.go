package tests

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/BZ6/golang/src/game"
	"github.com/BZ6/golang/src/web"
	"github.com/gorilla/websocket"
)

func Test1(t *testing.T) {
	server := CreateGameServer()
	defer server.Close()
	conn, _, err := CreateConnection(server)
	if err != nil {
		Fatal(t, "ошибка подключения к игре")
	}
	defer conn.Close()

	SetReadDuration(conn, 2 * time.Second)
	_, message, err := conn.ReadMessage()
	if err != nil {
		Fatal(t, "не удалось получить данные")
	}
	if len(message) == 0 {
		Fatal(t, "пустые данные")
	}

	output := string(message)
	if !strings.Contains(output, game.ALIVE) && !strings.Contains(output, game.DEAD) {
		Fatal(t, "данные без игровых символов")
	}

	Success(t, "web передал данные игры")
}

func Test2(t *testing.T) {
	server := CreateGameServer()
	defer server.Close()
	conn, _, err := CreateConnection(server)
	if err != nil {
		Fatal(t, "ошибка подключения к игре")
	}
	defer conn.Close()

	SetReadDuration(conn, 2 * time.Second)
	_, message, err := conn.ReadMessage()
	if err != nil {
		Fatal(t, "не удалось получить данные")
	}

	state := string(message)
	if len(state) == 0 {
		Fatal(t, "получено пустое состояние")
	}
	if !strings.Contains(state, game.ALIVE) {
		Fatal(t, "состояние не содержит живых клеток")
	}

	expectedMinSize := game.WIDTH * game.HEIGHT / 2
	if len(state) < expectedMinSize {
		Fatal(t, "состояние слишком маленькое")
	}

	Success(t, "web полностью передал состояние игры")
}

func Test3(t *testing.T) {
	server := CreateGameServer()
	defer server.Close()
	conn, _, err := CreateConnection(server)
	if err != nil {
		Fatal(t, "ошибка подключения к игре")
	}
	defer conn.Close()

	var states []string
	SetReadDuration(conn, 4 * time.Second)
	for i := 0; i < 2; i++ {
		_, message, err := conn.ReadMessage()
		if err != nil {
			t.Logf("Предупреждение: ошибка чтения состояния %d: %v", i, err)
			break
		}
		states = append(states, string(message))
	}

	if len(states) < 2 {
		Fatal(t, "получено недостаточно состояний")
	}
	if states[0] == states[1] {
		Fatal(t, "состояние не изменяется между кадрами")
	}

	Success(t, "web передает сменяющиеся состояния игры")
}

func Test4(t *testing.T) {
	var uni game.Universe
	uni.Init(10, 5)

	output := uni.ShowHTML()
	if output == "" {
		Fatal(t, "пустая строка в данных")
	}
	if !strings.Contains(output, game.ALIVE) || !strings.Contains(strings.ReplaceAll(output, game.ALIVE, ""), game.DEAD) {
		Fatal(t, "данные без игровых символов")
	}

	initialOutput := output
	uni.NextStep()
	newOutput := uni.ShowHTML()
	if initialOutput == newOutput {
		Fatal(t, "состояние не изменилось после шага")
	}

	Success(t, "игра верно отображается в HTML")
}

func Test5(t *testing.T) {
	server := CreateGameServer()
	defer server.Close()

	var connections []*websocket.Conn
	const clientCount = 3

	for i := 0; i < clientCount; i++ {
		conn, _, err := CreateConnection(server)
		if err != nil {
			Fatal(t, "подключение клиента к игре не удалось")
		}
		defer conn.Close()
		connections = append(connections, conn)
	}

	successfulClients := 0
	for _, conn := range connections {
		SetReadDuration(conn, 2 * time.Second)
		_, message, err := conn.ReadMessage()
		if err != nil {
			Fatal(t, "клиент не получил данные")
		}

		state := string(message)
		if len(state) == 0 {
			Fatal(t, "клиент получил пустые данные")
		}
		if !strings.Contains(state, game.ALIVE) && !strings.Contains(state, game.DEAD) {
			Fatal(t, "клиент получил данные без игровых символов")
		}

		successfulClients++
	}

	if successfulClients < clientCount {
		Fatal(t, "не все клиенты получили данные")
	}

	Success(t, "все клиенты получили данные от игры")
}

func Fatal(t *testing.T, msg string) {
	t.Fatalf("Интеграция нарушена: %s", msg)
}

func Success(t *testing.T, msg string) {
	t.Logf("Интеграция работает: %s", msg)
}

func CreateGameServer() (*httptest.Server) {
	return httptest.NewServer(http.HandlerFunc(web.GamePage))
}

func CreateConnection(server *httptest.Server) (*websocket.Conn, *http.Response, error) {
	wsURL := "ws" + server.URL[4:] + "/game"
	return websocket.DefaultDialer.Dial(wsURL, nil)
}

func SetReadDuration(conn *websocket.Conn, duration time.Duration) {
	conn.SetReadDeadline(time.Now().Add(duration))
}
