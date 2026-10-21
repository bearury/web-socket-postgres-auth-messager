package wshub

import (
	"slices"
	"web-socket-postgres-auth-messager/entities"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Client struct {
	Conn *websocket.Conn
	Send chan entities.WSMessageResponse
	User entities.UserResponse
}

type BroadcastExceptMessage struct {
	Message entities.WSMessageResponse
	Except  *Client
}

type SendToUserMessage struct {
	UserID  string
	Message entities.WSMessageResponse
}

type SendToChatMessage struct {
	UserIDs []string
	Message entities.WSMessageResponse
}

type Hub struct {
	clients         map[*Client]bool
	register        chan *Client
	unregister      chan *Client
	broadcast       chan entities.WSMessageResponse
	broadcastExcept chan BroadcastExceptMessage
	sendToUser      chan SendToUserMessage
	sendToChat      chan SendToChatMessage
}

func NewHub() *Hub {
	return &Hub{
		clients:         make(map[*Client]bool),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
		broadcast:       make(chan entities.WSMessageResponse),
		broadcastExcept: make(chan BroadcastExceptMessage),
		sendToUser:      make(chan SendToUserMessage),
		sendToChat:      make(chan SendToChatMessage),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

			logrus.Infof(
				"Client connected: userID=%d, clients=%d",
				client.User,
				len(h.clients),
			)

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)

				logrus.Infof(
					"Client disconnected: userID=%d, clients=%d",
					client.User,
					len(h.clients),
				)
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					logrus.Warnf(
						"Send buffer full: userID=%d",
						client.User,
					)
				}
			}

		case data := <-h.broadcastExcept:
			for client := range h.clients {
				if client == data.Except {
					continue
				}

				h.send(client, data.Message)
			}

		case data := <-h.sendToUser:
			for client := range h.clients {
				if client.User.Id == data.UserID {
					h.send(client, data.Message)
				}
			}

		case data := <-h.sendToChat:
			for client := range h.clients {
				// Проверкак что пользователь содержется в массиве переданных ID пользователей
				if slices.Contains(data.UserIDs, client.User.Id) {
					h.send(client, data.Message)
				}
			}
		}
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

func (h *Hub) Broadcast(message entities.WSMessageResponse) {
	h.broadcast <- message
}

func (client *Client) WritePump() {
	defer func(Conn *websocket.Conn) {
		err := Conn.Close()
		if err != nil {

		}
	}(client.Conn)

	for message := range client.Send {
		if err := client.Conn.WriteJSON(message); err != nil {
			logrus.Error("Error writing websocket message: ", err)
			return
		}
	}
}

func (h *Hub) send(client *Client, message entities.WSMessageResponse) {
	select {
	case client.Send <- message:
	default:
		logrus.Warnf(
			"Send buffer full: userID=%d",
			client.User,
		)
	}
}

func (h *Hub) BroadcastExcept(
	client *Client,
	message entities.WSMessageResponse,
) {
	h.broadcastExcept <- BroadcastExceptMessage{
		Message: message,
		Except:  client,
	}
}

func (h *Hub) SendToUser(
	userID string,
	message entities.WSMessageResponse,
) {
	h.sendToUser <- SendToUserMessage{
		UserID:  userID,
		Message: message,
	}
}

func (h *Hub) SendToChat(
	userIDs []string,
	message entities.WSMessageResponse,
) {
	h.sendToChat <- SendToChatMessage{
		UserIDs: userIDs,
		Message: message,
	}
}
