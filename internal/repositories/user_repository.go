package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/AlexWilliam12/silent-signal/internal/database"
	"github.com/AlexWilliam12/silent-signal/internal/dtos"
	"github.com/AlexWilliam12/silent-signal/internal/models"
	"github.com/lib/pq"
)

func CreateUser(user *models.User) (int64, error) {
	db := database.OpenConn()
	result := db.Select("username", "password", "credentials_hash").Create(user)
	return result.RowsAffected, result.Error
}

func FetchUserData(username string) (*dtos.UserResponse, error) {
	db, err := sql.Open("postgres", fmt.Sprintf(`
	user=%s
	dbname=%s
	password=%s
	sslmode=disable`,
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD")))
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
	SELECT
		username,
		credentials_hash,
		picture,
		ARRAY(
			SELECT JSON_BUILD_OBJECT(
				'contact',
				c.username,
				'picture',
				c.picture
			)
			FROM (
				SELECT 
					c.username, 
					c.picture
				FROM users c
				JOIN user_contacts uc ON uc.user_id = c.id
				WHERE uc.user_id = u.id
			) AS c
		) AS contacts,
		ARRAY(
			SELECT JSON_BUILD_OBJECT(
				'type',
				m.type,
				'content',
				m.content,
				'createdAt',
				m.created_at,
				'sender',
				m.sender_username,
				'senderPicture',
				m.sender_picture,
				'recipient',
				m.receiver_username,
				'recipientPicture',
				m.receiver_picture
			)
			FROM (
				SELECT
					pv.type,
					pv.content,
					pv.created_at,
					s.username AS sender_username,
    				s.picture AS sender_picture,
    				r.username AS receiver_username,
    				r.picture AS receiver_picture
				FROM private_messages pv
				LEFT JOIN users s ON pv.sender_id = s.id
				LEFT JOIN users r ON pv.receiver_id = r.id
				WHERE u.id = pv.sender_id OR u.id = pv.receiver_id AND pv.is_pending = FALSE
			) AS m
		) AS messages
	FROM users u
	WHERE u.username = $1`
	var (
		user            string
		credentialsHash string
		picture         sql.NullString
		contacts        pq.StringArray
		messages        pq.StringArray
	)
	err = db.QueryRow(query, username).Scan(
		&user,
		&credentialsHash,
		&picture,
		&contacts,
		&messages,
	)
	if err != nil {
		return nil, err
	}
	contactsResponse := make([]dtos.ContactResponse, 0)
	for _, contact := range contacts {
		var c dtos.ContactResponse
		if err = json.Unmarshal([]byte(contact), &c); err != nil {
			return nil, err
		}
		contactsResponse = append(contactsResponse, c)
	}

	messagesResponse := make([]dtos.PrivateMessageResponse, 0)
	for _, message := range messages {
		var p dtos.PrivateMessageResponse
		if err = json.Unmarshal([]byte(message), &p); err != nil {
			return nil, err
		}
		messagesResponse = append(messagesResponse, p)
	}

	return &dtos.UserResponse{
		Username:        user,
		CredentialsHash: credentialsHash,
		Picture:         picture.String,
		Contacts:        contactsResponse,
		Messages:        messagesResponse,
	}, nil
}

func FetchAllByUsernames(usernames []string) ([]*models.User, error) {
	db := database.OpenConn()
	var users []*models.User
	result := db.Where("username IN ?", usernames).Find(users)
	return users, result.Error
}

func FindUserByHash(request *dtos.CredentialsHashRequest) (*models.User, error) {
	db := database.OpenConn()
	var user models.User
	result := db.Where("credentials_hash = ?", request.Hash).Find(&user)
	return &user, result.Error
}

func FindUserByCredentials(request *dtos.UserRequest) (*models.User, error) {
	db := database.OpenConn()
	var user models.User
	result := db.Where("username = ? AND password = ?", request.Username, request.Password).First(&user)
	return &user, result.Error
}

func FindUserByName(username string) (*models.User, error) {
	db := database.OpenConn()
	var user models.User
	result := db.Where("username = ?", username).First(&user)
	return &user, result.Error
}

func UpdateUser(user *models.User) (int64, error) {
	db := database.OpenConn()
	result := db.Save(user)
	return result.RowsAffected, result.Error
}

func DeleteUserByName(username string) (int64, error) {
	db := database.OpenConn()
	result := db.Unscoped().Where("username = ?", username).Delete(&models.User{})
	return result.RowsAffected, result.Error
}

func SaveContact(user *models.User, contact *models.User) (int64, error) {
	db := database.OpenConn()
	user.Contacts = append(user.Contacts, contact)
	result := db.Save(user)
	return result.RowsAffected, result.Error
}
