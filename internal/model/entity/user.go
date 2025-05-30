package entity

type User struct {
}

func (user *User) TableName() string {
	return "users"
}