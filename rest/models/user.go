package models

type Users struct {
	ID        uint       `json:"id", gorm:"primary_key"`
	Email     string     `json:"email"`
	Firstname string     `json:"firstname"`
	Lastname  string     `json:"lastname"`
	Articles  []Articles `gorm:"foreignKey:UserId" json:"articles"`
}
