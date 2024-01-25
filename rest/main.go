package main

import (
	"rest/models"
	"rest/request/articles"
	"rest/request/users"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	models.ConnectDatabase()

	//ARTICLES
	articles.GetFiveArticles(r)
	articles.AddArticle(r)
	articles.GetAllArticles(r)
	articles.AddCommentToArticle(r)
	articles.UpdateArticleLikes(r)
	articles.UpdateArticleDislikes(r)
	articles.DeleteArticle(r)
	//USERS
	users.AddUser(r)
	users.Login(r)
	users.GetUserWithArticles(r)

	r.Run(":3000")

}
