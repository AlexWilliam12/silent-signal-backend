package main

import (
	"fmt"
	"net/http"

	"github.com/AlexWilliam12/silent-signal/internal/dtos"
	"github.com/gorilla/websocket"
)

type Response struct {
	Error   *string              `json:"error,omitempty"`
	Message *dtos.PrivateMessage `json:"message,omitempty"`
}

func main() {

	var token string
	fmt.Print("Your token: ")
	if _, err := fmt.Scanln(&token); err != nil {
		panic(err)
	}

	var sender string
	fmt.Print("Your name: ")
	if _, err := fmt.Scanln(&sender); err != nil {
		panic(err)
	}

	var recipient string
	fmt.Print("Recipient name: ")
	if _, err := fmt.Scanln(&recipient); err != nil {
		panic(err)
	}

	url := "ws://localhost:8080/chat/private"

	headers := http.Header{}
	headers.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(url, headers)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	go func() {
		for {
			var response Response
			err := conn.ReadJSON(&response)
			if err != nil {
				panic(err)
			}
			_, content, err := conn.ReadMessage()
			if err != nil {
				panic(err)
			}

			fmt.Println(string(content))
		}
	}()

	for {
		var input string
		_, err := fmt.Scanln(&input)

		if err != nil {
			panic(err)
		}

		if input == "/exit" {
			break
		}

		err = conn.WriteJSON(dtos.PrivateMessage{
			Sender:    sender,
			Recipient: recipient,
			Message: dtos.Message{
				Type:    "text",
				Content: input,
			},
		})
		if err != nil {
			panic(err)
		}

	}
}
