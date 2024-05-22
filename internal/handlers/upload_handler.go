package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/AlexWilliam12/silent-signal/internal/dtos"
	"github.com/AlexWilliam12/silent-signal/internal/models"
	"github.com/AlexWilliam12/silent-signal/internal/repositories"
	"github.com/AlexWilliam12/silent-signal/internal/services"
)

func HandleFetchPicture(w http.ResponseWriter, r *http.Request) {

	claims, err := HandleAuthorization(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	queryParams := r.URL.Query()

	groupParam := queryParams.Get("group")

	if groupParam != "" {

		group, err := repositories.FindGroupByName(groupParam)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if claims.Username != group.Creator.Username {
			http.Error(w, "unauthorized request", http.StatusForbidden)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&dtos.GroupResponse{Name: group.Name, Picture: group.Picture})

	} else {

		user, err := repositories.FindUserByName(claims.Username)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(&dtos.UserResponse{Username: user.Username, Picture: user.Picture})
	}
}

func HandleUploadPicture(w http.ResponseWriter, r *http.Request) {

	claims, err := HandleAuthorization(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	queryParams := r.URL.Query()

	groupParam := queryParams.Get("group")

	if groupParam != "" {

		group, err := repositories.FindGroupByName(groupParam)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if claims.Username != group.Creator.Username {
			http.Error(w, "unauthorized request", http.StatusForbidden)
			return
		}

		if group.Picture != "" {
			if err = services.DeleteFile(group.Picture[strings.LastIndex(group.Picture, "/"):]); err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

	} else {

		user, err := repositories.FindUserByName(claims.Username)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if user.Picture != "" {
			if err = services.DeleteFile(user.Picture[strings.LastIndex(user.Picture, "/"):]); err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

	}

	// 10MB max image size
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileType := fileHeader.Header.Get("Content-Type")
	if fileType != "image/jpg" && fileType != "image/jpeg" && fileType != "image/png" && fileType != "image/gif" {
		http.Error(w, "invalid image type, only those are permitted: jpg, jpeg, png, gif", http.StatusBadRequest)
		return
	}

	content, err := services.SaveFile(fileHeader)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if groupParam != "" {
		if _, err := repositories.SaveGroupPicture(groupParam, content); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		if _, err := repositories.SaveUserPicture(claims.Username, content); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func HandleChatUpload(w http.ResponseWriter, r *http.Request) {

	claims, err := HandleAuthorization(r)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	queryParams := r.URL.Query()

	groupParam := queryParams.Get("group")

	receiverParam := queryParams.Get("recipient")

	if receiverParam == "" && groupParam == "" {
		http.Error(w, "No recipient or group name specified", http.StatusBadRequest)
		return
	}

	if receiverParam != "" && groupParam != "" {
		http.Error(w, "Only one parameter can be specified at a time, whether recipient or group", http.StatusBadRequest)
		return
	}

	var recipient *models.User
	var group *models.Group

	if groupParam != "" {
		fetchedGroup, err := repositories.FindGroupByName(groupParam)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		group = fetchedGroup
	} else {
		user, err := repositories.FindUserByName(receiverParam)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		recipient = user
	}

	// 100MB max file size
	if err := r.ParseMultipartForm(100 << 20); err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	content, err := services.SaveFile(fileHeader)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fileType := fileHeader.Header.Get("Content-Type")

	if groupParam != "" {
		sender, err := repositories.FindUserByName(claims.Username)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = repositories.SaveGroupMessage(&models.GroupMessage{
			Sender:  *sender,
			Group:   *group,
			Type:    fileType,
			Content: content,
		})
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		message := dtos.GroupMessage{
			Sender: sender.Username,
			Group:  group.Name,
			Message: dtos.Message{
				Type:    fileType,
				Content: content,
			},
		}

		SendUploadMessage(groupParam, claims.Username, &message)

	} else {
		sender, err := repositories.FindUserByName(claims.Username)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = repositories.SavePrivateMessage(&models.PrivateMessage{
			Sender:    *sender,
			Recipient: *recipient,
			Type:      fileType,
			Content:   content,
		})
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
