package users

import (
	"net/http"
	"rest/dto"
	"rest/models"

	"github.com/gin-gonic/gin"
)

func GetUserByEmail(email string) (uint, error) {
	var user models.Users
	err := models.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

func AddUser(r *gin.Engine) {
	r.POST("/user", func(ctx *gin.Context) {
		var userDto dto.UserDto
		err := ctx.ShouldBindJSON(&userDto)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user := models.Users{
			Email:     userDto.Email,
			Firstname: userDto.Firstname,
			Lastname:  userDto.Lastname,
		}

		models.DB.Create(&user)

		ctx.JSON(http.StatusOK, gin.H{"Utilisateur sauvegardé ": user})
	})
}

func Login(r *gin.Engine) {
	r.POST("/login", func(ctx *gin.Context) {
		var user models.Users
		email := ctx.Query("email")
		if email == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Utilisateur connecté requis !"})
			return
		}
		//sauvegarde en bdd
		if err := models.DB.Where("email = ?", email).First(&user).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Aucun user existant !"})
			return
		}
		models.DB.Find(&user)
		ctx.JSON(http.StatusOK, gin.H{"Connexion ok ": "User connecté !"})

	})
}

func GetUserWithArticles(r *gin.Engine) {
	r.GET("/user/:userid/profile", func(ctx *gin.Context) {

		userId := ctx.Param("userid")

		if userId == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "ID utilisateur requis"})
			return
		}

		//récup infos user
		var user models.Users
		if err := models.DB.Where("id = ?", userId).First(&user).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur non trouvé"})
			return
		}

		//récup article de l'utilisateur
		var articles []models.Articles
		if err := models.DB.Where("user_id = ?", user.ID).Find(&articles).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur à la récupération des articles", "details": err.Error()})
			return
		}

		//affiche toutes les infos
		responseData := gin.H{
			"user":     user.Firstname,
			"articles": articles,
		}

		ctx.JSON(http.StatusOK, gin.H{"data": responseData})
	})
}
