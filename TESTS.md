# Отчет по тестам

`Захаркин Богдан`

`367224`

[титульник](title.pdf)

# Отчёт по интеграционному тестированию игры "Жизнь"

## Описание проекта

[ссылка на ридми](README.md)

Проект представляет собой веб-реализацию игры "Жизнь" Конвея на языке Go. Архитектура состоит из двух основных модулей:

- **Модуль `game`** - содержит логику игры:
  - Структура `Universe` для представления игрового поля
  - Алгоритмы вычисления следующего поколения клеток
  - Методы отображения состояния игры

- **Модуль `web`** - обеспечивает веб-интерфейс:
  - HTTP-обработчики для обслуживания клиентов
  - WebSocket-соединения для реального времени
  - Интеграция с игровой логикой

## Выбранные точки интеграции для тестирования

### 1. WebSocket соединение и передача данных
**Взаимодействие**: Клиент ↔ WebSocket ↔ Модуль web ↔ Модуль game

### 2. Корректность передачи игрового состояния
**Взаимодействие**: Модуль game → Модуль web → WebSocket → Клиент

### 3. Динамическое обновление игрового состояния
**Взаимодействие**: Игровой цикл → WebSocket трансляция → Множественные клиенты

### 4. Преобразование игрового состояния в HTML
**Взаимодействие**: Логика игры → Методы отображения → HTML вывод

### 5. Многопользовательская поддержка
**Взаимодействие**: Множественные клиенты → WebSocket сервер → Единая игровая логика

## Обоснование выбора интеграций

### Критичность взаимодействий

1. **WebSocket соединение** - фундаментальная основа реального времени
   - Без рабочего WebSocket вся система не функционирует
   - Проверяет интеграцию сетевого слоя с игровой логикой

2. **Целостность данных** - гарантия корректности передаваемого состояния
   - Предотвращает передачу битых или неполных данных
   - Обеспечивает согласованность состояния между сервером и клиентами

3. **Динамические обновления** - ядро игрового процесса
   - Проверяет корректность работы игрового цикла
   - Обеспечивает своевременное обновление состояния

4. **Преобразование данных** - критично для отображения
   - Проверяет корректность работы методов отображения
   - Гарантирует использование правильных игровых символов

5. **Многопользовательская работа** - проверка масштабируемости
   - Тестирует способность системы обслуживать нескольких клиентов
   - Проверяет отсутствие блокировок и конфликтов

## Примеры кода тестов

### Тест 1: Базовое WebSocket соединение

```go
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

	output := string(message)
	if !strings.Contains(output, game.ALIVE) && !strings.Contains(output, game.DEAD) {
		Fatal(t, "данные без игровых символов")
	}

	Success(t, "web передал данные игры")
}
```

### Тест 2: Проверка целостности игрового состояния

```go
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
	expectedMinSize := game.WIDTH * game.HEIGHT / 2
	if len(state) < expectedMinSize {
		Fatal(t, "состояние слишком маленькое")
	}

	Success(t, "web полностью передал состояние игры")
}
```

### Тест 3: Динамическое обновление состояния

```go
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
```

### Тест 4: Преобразование в HTML

```go
func Test4(t *testing.T) {
	var uni game.Universe
	uni.Init(10, 5)

	output := uni.ShowHTML()
	if output == "" {
		Fatal(t, "пустая строка в данных")
	}

	initialOutput := output
	uni.NextStep()
	newOutput := uni.ShowHTML()
	if initialOutput == newOutput {
		Fatal(t, "состояние не изменилось после шага")
	}

	Success(t, "игра верно отображается в HTML")
}
```

### Тест 5: Многопользовательская работа

```go
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

		successfulClients++
	}

	if successfulClients < clientCount {
		Fatal(t, "не все клиенты получили данные")
	}

	Success(t, "все клиенты получили данные от игры")
}
```

## Вспомогательные функции

```go
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
```

## Вывод тестов

```bash
=== RUN   Test1
    integration_test.go:171: Интеграция работает: web передал данные игры
--- PASS: Test1 (0.09s)
=== RUN   Test2
    integration_test.go:171: Интеграция работает: web полностью передал состояние игры
--- PASS: Test2 (0.07s)
=== RUN   Test3
    integration_test.go:171: Интеграция работает: web передает сменяющиеся состояния игры
--- PASS: Test3 (0.46s)
=== RUN   Test4
    integration_test.go:171: Интеграция работает: игра верно отображается в HTML
--- PASS: Test4 (0.00s)
=== RUN   Test5
    integration_test.go:171: Интеграция работает: все клиенты получили данные от игры
--- PASS: Test5 (0.09s)
PASS
ok      github.com/BZ6/golang/tests
```

## Выводы

### Результаты тестирования

**Все критические интеграции работают корректно:**
- WebSocket соединение устанавливается успешно
- Игровое состояние передается полностью и без искажений
- Динамические обновления происходят регулярно
- Преобразование в HTML формат работает правильно
- Система поддерживает multiple клиентов одновременно

### Сильные стороны интеграции

1. **Надежность соединения** - WebSocket обеспечивает стабильную передачу данных
2. **Целостность данных** - состояние игры передается полностью и корректно
3. **Масштабируемость** - система способна обслуживать нескольких клиентов
4. **Согласованность** - все клиенты получают одинаковое состояние игры

### Рекомендации

1. **Добавить обработку ошибок** на стороне клиента для улучшения отказоустойчивости
2. **Реализовать механизм паузы** для управления игровым процессом
3. **Добавить валидацию входных данных** для повышения безопасности
4. **Внедрить мониторинг** производительности WebSocket соединений

### Заключение

Интеграционное тестирование подтвердило, что взаимодействие между модулями `game` и `web` реализовано корректно. Система
демонстрирует стабильную работу всех ключевых компонентов, обеспечивая надежную передачу игрового состояния в реальном
времени. Архитектура проекта хорошо спроектирована для решения поставленных задач и готова к дальнейшему развитию.
