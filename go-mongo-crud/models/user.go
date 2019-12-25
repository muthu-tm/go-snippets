package models

type User struct {
	Name  string `json:"name"`
	ID    string `json:"id"`
	Email string `json:"emailID"`
}

func NewUser(name, id, email string) (newUser User) {
	newUser = User{name, id, email}
	return
}
