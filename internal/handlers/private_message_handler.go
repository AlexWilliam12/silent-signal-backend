package handlers

import (
	"log"
	"net/http"

	"github.com/AlexWilliam12/silent-signal/internal/dtos"
	"github.com/AlexWilliam12/silent-signal/internal/models"
	"github.com/AlexWilliam12/silent-signal/internal/repositories"
	"github.com/gorilla/websocket"
)

var (
	private_upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	privateConns     = make(map[string]*websocket.Conn)
	privateBroadcast = make(chan dtos.PrivateMessage)
)

func HandlePrivateConnections(w http.ResponseWriter, r *http.Request) {

	claims, err := HandleAuthorization(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	conn, err := private_upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	if _, ok := privateConns[claims.Username]; !ok {
		privateConns[claims.Username] = conn
	}

	go sendPendingPrivateMessages(claims.Username)

	for {
		var message dtos.PrivateMessage
		if err := conn.ReadJSON(&message); err != nil {
			log.Println(err)
			conn.Close()
			delete(privateConns, claims.Username)
			break
		}

		privateBroadcast <- message
	}
}

func HandlePrivateMessages() {
	for {
		message := <-privateBroadcast
		if conn, ok := privateConns[message.Recipient]; ok {
			if err := conn.WriteJSON(&message); err != nil {
				log.Println(err)
				conn.Close()
				delete(privateConns, message.Recipient)
				go savePrivateMessages(&message, true)
				break
			}
			go savePrivateMessages(&message, false)
		} else {
			go savePrivateMessages(&message, true)
		}
	}
}

func savePrivateMessages(message *dtos.PrivateMessage, isPending bool) {
	sender, err := repositories.FindUserByName(message.Sender)
	if err != nil {
		log.Println(err)
		return
	}

	receiver, err := repositories.FindUserByName(message.Recipient)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = repositories.SavePrivateMessage(&models.PrivateMessage{
		Sender:    *sender,
		Recipient: *receiver,
		Type:      "text",
		Content:   message.Message.Content,
		IsPending: isPending,
	})
	if err != nil {
		log.Println(err)
		return
	}
}

func sendPendingPrivateMessages(username string) {
	messages, err := repositories.FetchPendingPrivateMessages(username)
	if err != nil {
		log.Println(err)
		return
	}

	if len(messages) == 0 {
		return
	}

	if conn, ok := privateConns[username]; ok {
		var ids []uint
		for _, message := range messages {
			response := dtos.PrivateMessage{
				Sender:    message.Sender.Username,
				Recipient: message.Recipient.Username,
				Message: dtos.Message{
					Type:    "text",
					Content: message.Content,
				},
			}
			if err = conn.WriteJSON(&response); err != nil {
				log.Println(err)
				continue
			}
			ids = append(ids, message.ID)
		}
		go repositories.UpdatePendingSituation(ids)
	}
}
