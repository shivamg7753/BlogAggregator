package handlers

import (
	"blogAggregator/internal/auth"
	"blogAggregator/internal/database"
	"blogAggregator/internal/models"
	"blogAggregator/internal/rss"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}
	err := c.ShouldBind(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: hashedPassword,
	}
	err = database.DB.Create(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User already exists or invalid data",
		})
		return
	}
	
	// Don't return the password hash
	user.Password = ""
	c.JSON(http.StatusCreated, user)
}

//add feed

func CreateFeed(c *gin.Context) {
	var input struct {
		Title string `json:"title" binding:"required"`
		URL   string `json:"url" binding:"required"`
	}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	feed := models.Feed{Title: input.Title, URL: input.URL}
	err = database.DB.Create(&feed).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "feed already exists",
		})
		return
	}
	c.JSON(http.StatusCreated, feed)
}

//list feed

func ListFeeds(c *gin.Context) {
	var feeds []models.Feed
	err := database.DB.Find(&feeds).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, feeds)
}

func RefreshFeed(c *gin.Context) {
	var input struct {
		FeedId uint `json:"feed_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var feed models.Feed
	if err := database.DB.First(&feed, input.FeedId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	if err := rss.FetchAndStoreFeed(feed); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "feed refreshed",
	})
}

func ListPosts(c *gin.Context) {
	var posts []models.Post
	database.DB.Order("published desc").Limit(20).Find(&posts)
	c.JSON(http.StatusOK, posts)
}

func CreateUser(c *gin.Context) {
	var input struct {
		UserName string `json:"username" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	user := models.User{Username: input.UserName}
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func SubscribeFeed(c *gin.Context) {
	var input struct {
		UserID uint `json:"user_id" binding:"required"`
		FeedID uint `json:"feed_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	sub := models.Subscription{
		UserID: input.UserID,
		FeedID: input.FeedID,
	}

	if err := database.DB.Create(&sub).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, sub)
}

func GetUserFeed(c *gin.Context) {
	userId := c.GetUint("User_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	var FeedIDs []int
	database.DB.Model(&models.Subscription{}).Where("user_id=?", userId).Pluck("feed_id", &FeedIDs)

	var posts []models.Post
	 result:=database.DB.Where("feed_id IN ?", FeedIDs).
		Order("published desc").
		Limit(20).
		Offset(offset).
		Find(&posts)

   if result.Error!=nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	 }

	c.JSON(http.StatusOK, gin.H{
		"page":  page,
		"limit": limit,
		"posts": posts,
	})
}

func Login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Check password
	if !auth.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}
