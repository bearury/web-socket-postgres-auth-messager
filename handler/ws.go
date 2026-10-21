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
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logrus.Error("Error connection websocket: ", err)
		return
	}

	token := c.Query("token")

	userID, err := handler.services.Authorization.ParseToken(token)
	if err != nil {
		logrus.Error("Error parsing token: ", err)
		err := conn.Close()
		if err != nil {
			return
		}
		return
	}

	user, err := handler.services.GetUserById(userID)
	if err != nil {
		logrus.Error("Error get User by ID: ", err)
		err := conn.Close()
		if err != nil {
			return
		}
		return
	}

	client := &ws_hub.Client{
		Conn: conn,
		Send: make(chan entities.MessageResponse, 256),
		User: user,
	}

	handler.hub.Register(client)

	defer func() {
		handler.hub.Unregister(client)
	}()

	go client.WritePump()

	if err := handler.readFromClient(client); err != nil {
		logrus.Error("Error reading from websocket: ", err)
	}

	logrus.Infof(
		"Websocket connection closed: %s",
		c.Request.RemoteAddr,
	)
}

func (handler *Handler) readFromClient(client *ws_hub.Client) error {
	for {
		msg := new(entities.Message)

		logger.Info("UserID: %s", client.User.Id)

		if err := client.Conn.ReadJSON(msg); err != nil {
			return err
		}

		switch msg.Type {
		case "message":
			content, ok := msg.Payload["content"]
			if !ok {
				continue
			}

			logger.Info("CONTENT !!!!!!!!!!!!", content)

			chatId, ok := msg.Payload["chat_id"]
			if !ok {
				continue
			}

			contentStr, ok := content.(string)
			if !ok {
				logrus.Error("Message content is not string")
				continue
			}

			chatIdStr, ok := chatId.(string)
			if !ok {
				logrus.Error("Chat id is not string")
				continue
			}

			logger.Info("CONTENT ☀", contentStr)
			logger.Info("CHAT ☀", chatIdStr)

			handler.sendToChatMessage(chatIdStr, client, contentStr)

		case "ping":
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

func (handler *Handler) sendAllMessage(client *ws_hub.Client, content string) {
	name := "Anonimus"

	if client.User.Name != "" {
		name = client.User.Name
	} else if client.User.Username != "" {
		name = client.User.Username
	}

	handler.hub.SendToUser(client.User.Id,
		entities.MessageResponse{
			Type: "message",
			Payload: entities.PayloadMessageResponse{
				From: name,
				To: entities.PayloadToMessage{
					Type: entities.ChatTypeChat,
					Id:   client.User.Id,
				},
				Content: content,
			},
		})
}

func (handler *Handler) sendToChatMessage(chatID string, client *ws_hub.Client, content string) {

	users, err := handler.services.Chat.GetUsersChatByIdWithoutUser(chatID, client.User.Id)

	if err != nil {
		logrus.Error("Error Get Chat_users By Id Without current User. Error:", err)
	}

	userIDs := make([]string, len(users))

	for i, user := range users {
		userIDs[i] = user.UserId
	}

	handler.hub.SendToChat(userIDs,
		entities.MessageResponse{
			Type: "message",
			Payload: entities.PayloadMessageResponse{
				From: client.User.Username,
				To: entities.PayloadToMessage{
					Type: entities.ChatTypeChat,
					Id:   chatID,
				},
				Content: content,
			},
		})
}
