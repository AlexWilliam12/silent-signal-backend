package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AlexWilliam12/silent-signal/internal/dtos"
	"github.com/AlexWilliam12/silent-signal/internal/models"
	"github.com/AlexWilliam12/silent-signal/internal/repositories"
	"github.com/gorilla/websocket"
)

var (
	group_upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	groupConns = make(map[string][]*dtos.GroupUser)
)

func HandleGroupMessages(w http.ResponseWriter, r *http.Request) {

	conn, err := group_upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	claims, err := HandleAuthorization(r)
	if err != nil {
		log.Println(err)
		conn.WriteJSON(fmt.Sprintf(`{"error":"%s"}`, err.Error()))
		return
	}

	queryParams := r.URL.Query()
	groupParam := queryParams.Get("name")

	if groupParam == "" {
		conn.WriteJSON(`{"error":"group name not specified"}`)
		return
	}

	group, err := repositories.FindGroupByName(groupParam)
	if err != nil {
		log.Println(err)
		conn.WriteJSON(fmt.Sprintf(`{"error":"%s"}`, err.Error()))
		return
	}

	if conns, ok := groupConns[groupParam]; ok {
		var isPresent bool
		for _, groupConn := range conns {
			if groupConn.Conn == conn {
				isPresent = true
				break
			}
		}
		if !isPresent {
			groupConns[groupParam] = append(groupConns[groupParam], &dtos.GroupUser{Username: claims.Username, Conn: conn})
		}
	} else {
		groupConns[groupParam] = append(groupConns[groupParam], &dtos.GroupUser{Username: claims.Username, Conn: conn})
	}

	go sendPendingGroupMessages(claims.Username, conn)

	for {
		var message dtos.GroupMessage
		if err := conn.ReadJSON(&message); err != nil {
			log.Println(err)
			break
		}

		if conns, ok := groupConns[message.Group]; ok {
			for _, groupConn := range conns {
				if groupConn.Conn != conn {
					if err := groupConn.Conn.WriteJSON(&message); err != nil {
						log.Println(err)
						groupConn.Conn.Close()
						delete(groupConns, message.Group)
						break
					}
				}
			}
			go saveGroupMessages(&message, group)
		}
	}
}

func saveGroupMessages(message *dtos.GroupMessage, group *models.Group) {

	sender, err := repositories.FindUserByName(message.Sender)
	if err != nil {
		log.Println(err)
		return
	}

	var usernames []string
	for _, conn := range groupConns[group.Name] {
		usernames = append(usernames, conn.Username)
	}

	users, err := repositories.FetchAllByUsernames(usernames)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = repositories.SaveGroupMessage(&models.GroupMessage{
		Sender:  *sender,
		Group:   *group,
		Type:    "text",
		Content: message.Message.Content,
		SeenBy:  users,
	})
	if err != nil {
		log.Println(err)
		return
	}
}

func sendPendingGroupMessages(username string, conn *websocket.Conn) {

	messages, err := repositories.FetchPendingGroupMessages(username)
	if err != nil {
		log.Println(err)
		return
	}

	for _, message := range messages {
		err := conn.WriteJSON(dtos.GroupMessage{
			Sender: message.Sender.Username,
			Group:  message.Group.Name,
			Message: dtos.Message{
				Type:    message.Type,
				Content: message.Content,
			},
		})
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

func SendUploadMessage(groupName string, username string, message *dtos.GroupMessage) {
	if conns, ok := groupConns[groupName]; ok {
		for _, conn := range conns {
			if conn.Username != username {
				conn.Conn.WriteJSON(message)
			}
		}
	}
}
