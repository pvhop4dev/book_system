package entity

type Book struct {
}

func (book *Book) TableName() string {
	return "books"
}
