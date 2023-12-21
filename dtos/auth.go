package dtos

type Users struct {
	Username string `json:"email"`
}

type PasswordData struct {
	Password string `json:"password"`
}

type CodeData struct {
	ConfirmationCode string `json:"code"`
}

type AuthData struct {
	Users
	PasswordData
}

type AuthCodeData struct {
	Users
	CodeData
}

type AuthResetData struct {
	AuthData
	CodeData
}

type UpdateGroup struct {
	Groups string `json:"groups"`
}
