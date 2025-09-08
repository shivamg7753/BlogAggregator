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

// Request DTOs for Swagger
type RegisterInput struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

type FeedCreateInput struct {
    Title string `json:"title"`
    URL   string `json:"url"`
}

type RefreshFeedInput struct {
    FeedId uint `json:"feed_id"`
}

type CreateUserInput struct {
    UserName string `json:"username"`
}

type SubscribeInput struct {
    UserID uint `json:"user_id"`
    FeedID uint `json:"feed_id"`
}

type LoginInput struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

// RegisterUser
// @Summary      Register user
// @Description  Create a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input  body  RegisterInput  true  "User payload"
// @Success      201    {object}  models.User
// @Failure      400    {object}  map[string]string
// @Router       /users/register [post]
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

// CreateFeed
// @Summary      Create feed
// @Tags         feeds
// @Accept       json
// @Produce      json
// @Param        input  body  FeedCreateInput  true  "Feed"
// @Success      201    {object}  models.Feed
// @Failure      400    {object}  map[string]string
// @Router       /feeds [post]
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

// ListFeeds
// @Summary      List feeds
// @Tags         feeds
// @Produce      json
// @Success      200  {array}  models.Feed
// @Router       /feeds [get]
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

// RefreshFeed
// @Summary      Refresh a feed
// @Tags         feeds
// @Accept       json
// @Produce      json
// @Param        input body RefreshFeedInput true "Feed to refresh"
// @Success      200  {object} map[string]string
// @Failure      400  {object} map[string]string
// @Router       /feeds/refresh [post]
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

// ListPosts
// @Summary      List latest posts
// @Tags         posts
// @Produce      json
// @Success      200  {array}  models.Post
// @Router       /posts [get]
func ListPosts(c *gin.Context) {
	var posts []models.Post
	database.DB.Order("published desc").Limit(20).Find(&posts)
	c.JSON(http.StatusOK, posts)
}

// CreateUser
// @Summary      Create user (internal/demo)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        input body CreateUserInput true "User"
// @Success      201 {object} models.User
// @Router       /users [post]
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

// SubscribeFeed
// @Summary      Subscribe to a feed
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        input body SubscribeInput true "Subscription"
// @Success      201 {object} models.Subscription
// @Failure      401 {object} map[string]string
// @Router       /subscriptions [post]
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

// GetUserFeed
// @Summary      Get personalized feed
// @Tags         users
// @Produce      json
// @Param        id    path      int     true  "User ID"
// @Param        page  query     int     false "Page"
// @Param        limit query     int     false "Limit"
// @Success      200   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]string
// @Router       /users/{id}/feed [get]
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

// Login
// @Summary      User login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        input body LoginInput true "Credentials"
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]string
// @Router       /login [post]
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
