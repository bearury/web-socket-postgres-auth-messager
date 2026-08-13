package handler

import (
	"net/http"
	"web-socket-postgres-auth-messager/entities"
	ws_hub "web-socket-postgres-auth-messager/handler/ws-hub"

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
		conn.Close()
		return
	}

	user, err := handler.services.GetUserById(userID)
	if err != nil {
		logrus.Error("Error get User by ID: ", err)
		conn.Close()
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

		if err := client.Conn.ReadJSON(msg); err != nil {
			return err
		}

		switch msg.Type {
		case "message":
			payload, ok := msg.Payload["content"]
			if !ok {
				continue
			}

			content, ok := payload.(string)
			if !ok {
				logrus.Error("Message content is not string")
				continue
			}

			logrus.Infof(
				"User %d sent message: %s",
				client.User,
				content,
			)

			name := "Anonimus"

			if client.User.Name != "" {
				name = client.User.Name
			} else if client.User.Username != "" {
				name = client.User.Username
			}

			handler.hub.BroadcastExcept(client,
				entities.MessageResponse{
					Type: "message",
					Payload: entities.PayloadMessageResponse{
						From:    name,
						To:      "all",
						Content: content,
					},
				})

		case "ping":
			timestamp, ok := msg.Payload["timestamp"]
			if !ok {
				continue
			}

			logrus.Infof(
				"Ping from user %d: %v",
				client.User,
				timestamp,
			)
		}
	}
}
