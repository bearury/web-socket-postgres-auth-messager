package wshub

import (
	"web-socket-postgres-auth-messager/entities"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type Client struct {
	Conn *websocket.Conn
	Send chan entities.MessageResponse
	User entities.UserResponse
}

type BroadcastExceptMessage struct {
	Message entities.MessageResponse
	Except  *Client
}

type SendToUserMessage struct {
	UserID  string
	Message entities.MessageResponse
}

type Hub struct {
	clients         map[*Client]bool
	register        chan *Client
	unregister      chan *Client
	broadcast       chan entities.MessageResponse
	broadcastExcept chan BroadcastExceptMessage
	sendToUser      chan SendToUserMessage
}

func NewHub() *Hub {
	return &Hub{
		clients:         make(map[*Client]bool),
		register:        make(chan *Client),
		unregister:      make(chan *Client),
		broadcast:       make(chan entities.MessageResponse),
		broadcastExcept: make(chan BroadcastExceptMessage),
		sendToUser:      make(chan SendToUserMessage),
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
		}
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

func (h *Hub) Broadcast(message entities.MessageResponse) {
	h.broadcast <- message
}

func (client *Client) WritePump() {
	defer client.Conn.Close()

	for message := range client.Send {
		if err := client.Conn.WriteJSON(message); err != nil {
			logrus.Error("Error writing websocket message: ", err)
			return
		}
	}
}

func (h *Hub) send(client *Client, message entities.MessageResponse) {
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
	message entities.MessageResponse,
) {
	h.broadcastExcept <- BroadcastExceptMessage{
		Message: message,
		Except:  client,
	}
}

func (h *Hub) SendToUser(
	userID string,
	message entities.MessageResponse,
) {
	h.sendToUser <- SendToUserMessage{
		UserID:  userID,
		Message: message,
	}
}
