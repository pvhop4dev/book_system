package restapi

import (
	"book_system/internal/model"
	"book_system/internal/service"
	"book_system/internal/transport/middleware"
	"book_system/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService service.IUserService
}

func NewUserController(userService service.IUserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (uc *UserController) SetupUsersRoutes(router *gin.RouterGroup) {
	router.GET("/me", uc.GetUserProfile)
	router.PUT("/me", uc.UpdateUserProfile)
	router.GET("", uc.ListUsers)
}

func (uc *UserController) SetupAuthRoutes(router *gin.RouterGroup) {
	router.POST("/register", uc.Register)
	router.POST("/login", uc.Login)
	router.POST("/refresh", uc.RefreshToken)
}

// RegisterUser godoc
// @Summary Register a new user
// @Description Register a new user with the input payload
// @Tags users
// @Accept  json
// @Produce  json
// @Param input body model.RegisterRequest true "Register info"
// @Success 201 {object} model.RegisterResponse
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /api/v1/auth/register [post]
func (uc *UserController) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errT := middleware.BadRequest
		c.JSON(errT.Code, gin.H{"error": errT.GetMesssageI18n(utils.GetCurrentLang(c))})
		c.Abort()
		return
	}

	resp, err := uc.userService.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// Login godoc
// @Summary Login a user
// @Description Login with email and password
// @Tags auth
// @Accept  json
// @Produce  json
// @Param input body model.LoginRequest true "Login credentials"
// @Success 200 {object} model.LoginResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /api/v1/auth/login [post]
func (uc *UserController) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := uc.userService.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// GetUserProfile godoc
// @Summary Get user profile
// @Description Get the profile of the authenticated user
// @Tags users
// @Security BearerAuth
// @Produce  json
// @Success 200 {object} model.User
// @Failure 401 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Router /api/v1/users/me [get]
func (uc *UserController) GetUserProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := uc.userService.GetUserByID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUserProfile godoc
// @Summary Update user profile
// @Description Update the profile of the authenticated user
// @Tags users
// @Security BearerAuth
// @Accept  json
// @Produce  json
// @Param input body model.UpdateUserRequest true "User update info"
// @Success 200 {object} model.User
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /api/v1/users/me [put]
func (uc *UserController) UpdateUserProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req model.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	updatedUser, err := uc.userService.UpdateUser(c.Request.Context(), userID.(string), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// ListUsers godoc
// @Summary List all users
// @Description Get a list of users with pagination
// @Tags users
// @Security BearerAuth
// @Produce  json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /api/v1/users [get]
func (uc *UserController) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	users, total, err := uc.userService.ListUsers(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": users,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_page": (int(total) + pageSize - 1) / pageSize,
		},
	})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Refresh access token using refresh token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param input body model.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} model.LoginResponse
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Router /api/v1/auth/refresh [post]
func (uc *UserController) RefreshToken(c *gin.Context) {
	var req model.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := uc.userService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
