package model

type User struct {
	ID         uint
	Name       string
	Notes      []Note
	CreditCard CreditCard
}

type Note struct {
	ID      uint
	Title   string
	Content string
	UserID  uint
}

type CreditCard struct {
	ID     uint
	Number string
	UserID uint
}
