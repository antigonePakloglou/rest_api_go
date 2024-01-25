package models

type Comments struct {
	ID        uint   `json:"id", gorm:"primary_key"`
	Content   string `json:"content"`
	ArticleId uint   `json:"article_id"`
}
