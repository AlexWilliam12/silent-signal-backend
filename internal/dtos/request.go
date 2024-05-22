package dtos

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type GroupRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ContactRequest struct {
	Contact string `json:"contact"`
}

type RecoverPasswordRequest struct {
	Hash        string `json:"hash"`
	NewPassword string `json:"newPassword"`
}

type CredentialsHashRequest struct {
	Hash string `json:"hash"`
}
