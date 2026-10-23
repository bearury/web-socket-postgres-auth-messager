package handler

import (
	"net/http"
	"web-socket-postgres-auth-messager/entities"
	ws_hub "web-socket-postgres-auth-messager/handler/ws-hub"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func (handler *Handler) ws(c *gin.Context) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Для разработки ок, для продакшена лучше ограничить домены
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.Error("Error connection websocket: ", err)
		return
	}
	// Важно: гарантируем закрытие сокета при выходе из хэндлера,
	// если мы не передали управление в read/write циклы
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	token := c.Query("token")

	userID, err := handler.services.Authorization.ParseToken(token)
	if err != nil {
		logrus.Error("Error parsing token: ", err)
		// Вежливо говорим клиенту, что токен невалиден
		err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.ClosePolicyViolation, "Unauthorized"))
		if err != nil {
			return
		}
		return
	}

	user, err := handler.services.GetUserById(userID)
	if err != nil {
		logrus.Error("Error get User by ID: ", err)
		err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "User not found"))
		if err != nil {
			return
		}
		return
	}

	client := &ws_hub.Client{
		Conn: conn,
		Send: make(chan entities.WSMessage, 256),
		User: user,
	}

	// Регистрируем клиента в хабе
	handler.hub.Register(client)

	// Правильный порядок завершения:
	// 1. Чтение (readFromClient) завершается (клиент ушел или ошибка).
	// 2. Срабатывает defer: удаляем клиента из хаба.
	// 3. Внутри Unregister ОБЯЗАТЕЛЬНО должен закрываться канал client.Send
	//    ИЛИ само соединение conn, что послужит сигналом для выхода из WritePump.
	defer func() {
		handler.hub.Unregister(client)
		logrus.Infof("Websocket connection closed for user %d: %s", user.Id, c.Request.RemoteAddr)
	}()

	// Запускаем запись в фоне
	go client.WritePump()

	// Сам хэндлер блокируется на чтении. Это главный цикл.
	if err := handler.readFromClient(client); err != nil {
		// Не пишем в лог ошибку, если это обычное закрытие (1001 Going Away или 1000 Normal)
		if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
			logrus.Error("Error reading from websocket: ", err)
		}
	}
}

func (handler *Handler) readFromClient(client *ws_hub.Client) error {
	for {
		msg := new(entities.WSMessage)

		logger.Info("UserID: %s", client.User.Id)

		if err := client.Conn.ReadJSON(msg); err != nil {
			return err
		}

		switch msg.Type {
		case entities.WSTypeMessage:
			err := handler.services.WebSocket.HandleMessage(client, msg)
			if err != nil {
				logrus.Error("Error websocket handle message: ", err)
				return nil
			}
		case entities.WSTypeNotify:

		case entities.WSTypePing:
			//timestamp, ok := msg.Payload["timestamp"]
			//if !ok {
			//	continue
			//}

			//logrus.Infof(
			//	"Ping from user %d: %v",
			//	client.User,
			//	timestamp,
			//)
		}
	}
}

//func (handler *Handler) sendAllMessage(client *ws_hub.Client, content string) {
//	name := "Anonimus"
//
//	if client.User.Name != "" {
//		name = client.User.Name
//	} else if client.User.Username != "" {
//		name = client.User.Username
//	}
//
//	handler.hub.SendToUser(client.User.Id,
//		entities.WSMessage{
//			Type: "message",
//			Payload: {
//				From: name,
//				To: entities.WSPayloadToMessage{
//					Type: entities.ChatTypeChat,
//					Id:   client.User.Id,
//				},
//				Content: content,
//			},
//		})
//}
