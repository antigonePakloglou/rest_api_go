package articles

import (
	"net/http"
	"rest/dto"
	"rest/models"
	"rest/request/users"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetFiveArticles(r *gin.Engine) {
	r.GET("/home", func(ctx *gin.Context) {
		var articles []models.Articles

		// Retrieve the first 5 articles from the database
		if err := models.DB.Limit(5).Find(&articles).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error de recuperation d' articles", "details": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": articles})
	})
}

func AddArticle(r *gin.Engine) {
	r.POST("/article", func(ctx *gin.Context) {
		var articleDto dto.ArticleDto
		err := ctx.ShouldBindJSON(&articleDto)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//recupère id utilisateur
		//récupère email
		email := ctx.Query("email")
		userId, err := users.GetUserByEmail(email)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur innexistant"})
			return
		}
		articleDto.UserId = userId

		article := models.Articles{
			Content: articleDto.Content,
			Likes:   0,
			UserId:  articleDto.UserId,
		}

		models.DB.Create(&article)

		ctx.JSON(http.StatusOK, gin.H{"data": article})
	})
}

func GetAllArticles(r *gin.Engine) {
	r.GET("/articles", func(ctx *gin.Context) {
		var articles []models.Articles
		//récupère email
		email := ctx.Query("email")
		if email == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Utilisateur connecté requis !"})
			return
		}
		//recupère id utilisateur
		userId, err := users.GetUserByEmail(email)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Utilisateur innexistant"})
			return
		}
		if err := models.DB.Where("user_id = ?", userId).Find(&articles).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des articles", "details": err.Error()})
			return
		}
		//vient remplir le tableau d'articles
		models.DB.Find(&articles)
		ctx.JSON(http.StatusOK, gin.H{"Article ID": articles})
	})
}

func AddCommentToArticle(r *gin.Engine) {
	r.POST("/articles/:articleid/comment", func(ctx *gin.Context) {
		var commentDto dto.CommentDto
		articleID := ctx.Param("articleid")

		if articleID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Article ID requis"})
			return
		}

		// conversion vers uint64
		idArticle64, err := strconv.ParseUint(articleID, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Article ID invalide"})
			return
		}

		//puis vers 32
		idArticle32 := uint(idArticle64)
		commentDto.ArticleId = idArticle32

		//lien entre json et dto
		if err := ctx.ShouldBindJSON(&commentDto); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Erreur de mapping", "details": err.Error()})
			return
		}

		comment := models.Comments{
			Content:   commentDto.Content,
			ArticleId: commentDto.ArticleId,
		}

		// Création du commentaire dans la base de données
		if err := models.DB.Create(&comment).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la création du commentaire", "details": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": comment})
	})
}

func UpdateArticleLikes(r *gin.Engine) {
	r.POST("/articles/:articleid/like", func(ctx *gin.Context) {
		articleID := ctx.Param("articleid")

		if articleID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Article ID requis"})
			return
		}

		//conversion uint64
		idArticle64, err := strconv.ParseUint(articleID, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Article ID invalide"})
			return
		}

		//puis uint32
		idArticle32 := uint(idArticle64)

		// récupérer l'article
		var articleExistant models.Articles
		if err := models.DB.First(&articleExistant, idArticle32).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Article non trouvé"})
			return
		}

		//ajout d'un like
		articleExistant.Likes++

		//sauvegarde en bdd
		if err := models.DB.Save(&articleExistant).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur a l'update", "details": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": articleExistant})
	})
}

func UpdateArticleDislikes(r *gin.Engine) {
	r.POST("/articles/:articleid/dislike", func(ctx *gin.Context) {
		articleID := ctx.Param("articleid")

		if articleID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Article ID requis"})
			return
		}

		//conversion uint64
		idArticle64, err := strconv.ParseUint(articleID, 10, 32)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Article ID invalide"})
			return
		}

		//puis uint32
		idArticle32 := uint(idArticle64)

		// récupérer l'article
		var articleExistant models.Articles
		if err := models.DB.First(&articleExistant, idArticle32).Error; err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Article non trouvé"})
			return
		}

		//vérification si nombre de likes est supérieur à 0
		if articleExistant.Likes > 0 {
			articleExistant.Likes--

			//mise a jour article avec like en moin
			if err := models.DB.Save(&articleExistant).Error; err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur à la mise à jour des likes de l'article", "details": err.Error()})
				return
			}
		} else {
			//sinon erreur
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Le nombre de likes est déjà 0"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"data": articleExistant})
	})
}

func DeleteArticle(r *gin.Engine) {
	r.DELETE("/articles/:id", func(ctx *gin.Context) {
		//récupération de l'id de l'article
		articleID := ctx.Param("id")
		//verif si parametre ok
		if articleID == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Article ID nécessaire"})
			return
		}

		//converstion de l'id en uint
		id, err := strconv.ParseUint(articleID, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Article ID invalide"})
			return
		}

		//suppression des commentaires liés à l'article
		if err := models.DB.Where("article_id = ?", id).Delete(&models.Comments{}).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur à la suppression des commentaires", "details": err.Error()})
			return
		}

		//enfin suppression de l'article
		if err := models.DB.Where("id = ?", id).Delete(&models.Articles{}).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur à la suppression de l'article", "details": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Article supprimé correctement !"})
	})
}
