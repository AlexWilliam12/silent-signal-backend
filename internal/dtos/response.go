package dtos

import "time"

type UserResponse struct {
	Username        string                   `json:"username"`
	CredentialsHash string                   `json:"credentialsHash"`
	Picture         string                   `json:"picture"`
	Messages        []PrivateMessageResponse `json:"messages"`
	Contacts        []ContactResponse        `json:"contacts"`
}

type ContactResponse struct {
	Contact string `json:"contact"`
	Picture string `json:"picture"`
}

type PrivateMessageResponse struct {
	Sender           string    `json:"sender"`
	SenderPicture    string    `json:"senderPicture"`
	Recipient        string    `json:"recipient"`
	RecipientPicture string    `json:"recipientPicture"`
	Type             string    `json:"type"`
	Content          string    `json:"content"`
	CreatedAt        time.Time `json:"createdAt"`
}

type GroupResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Creator     string `json:"creator"`
	Picture     string `json:"picture"`
}
