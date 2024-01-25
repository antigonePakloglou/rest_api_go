package dto

type ArticleDto struct {
	Content string `json:"content"`
	Likes   uint   `json:"likes"`
	UserId  uint   `json:"user_id"`
}

type CommentDto struct {
	Content   string `json:"content"`
	ArticleId uint   `json:"article_id"`
}

type UserDto struct {
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}
