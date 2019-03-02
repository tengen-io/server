package models

type User struct {
	IdentityID int32
	Name       string `json:"name"`
}

type Identity struct {
	Id    int32  `json:"id"`
	Email string `json:"email"`
	User  `json:"user"`
}
