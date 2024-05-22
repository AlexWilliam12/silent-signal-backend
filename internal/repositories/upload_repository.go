package repositories

import (
	"github.com/AlexWilliam12/silent-signal/internal/database"
	"github.com/AlexWilliam12/silent-signal/internal/models"
)

func SaveUserPicture(username string, imgUrl string) (int64, error) {
	db := database.OpenConn()
	result := db.Model(&models.User{}).Where("username = ?", username).Update("picture", imgUrl)
	return result.RowsAffected, result.Error
}

func SaveGroupPicture(groupName string, imgUrl string) (int64, error) {
	db := database.OpenConn()
	result := db.Model(&models.Group{}).Where("name = ?", groupName).Update("picture", imgUrl)
	return result.RowsAffected, result.Error
}
