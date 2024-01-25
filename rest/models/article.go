package models

type Articles struct {
	ID       uint       `json:"id", gorm:"primary_key"`
	Content  string     `json:"content"`
	Likes    uint       `json:"likes"`
	UserId   uint       `json:"user_id"`
	Comments []Comments `gorm:"foreignKey:ArticleId" json:"comments"`
}
