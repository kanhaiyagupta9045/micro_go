package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/kanhaiyagupta9045/pratilipi/user/internal/helpers"
	"github.com/kanhaiyagupta9045/pratilipi/user/internal/model"
	"github.com/kanhaiyagupta9045/pratilipi/user/internal/service"
	"gorm.io/gorm"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

var validate = validator.New()

func (u *UserHandler) RegisterUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user model.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if validationErr := validate.Struct(user); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		if err := u.service.CreateUser(&user); err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{"erorr": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message:": "User Registered Successfully"})
	}
}

func (u *UserHandler) ListAllUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := u.service.GetAllUser()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	}
}

func (u *UserHandler) GetUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide id"})
			return
		}
		userID, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id format"})
			return
		}

		user, err := u.service.GetUserByID(uint(userID))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func (u *UserHandler) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var logindata model.LoginData

		if err := c.BindJSON(&logindata); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if validationErr := validate.Struct(logindata); validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		user, err := u.service.LoginUser(logindata)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		token, err := helpers.GenerateAccessToken(int(user.ID))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"accesstoken": token})

	}
}

func (u *UserHandler) ValidateToken() gin.HandlerFunc {
	return func(c *gin.Context) {

		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide token for validating the users"})
			return
		}

		claims, err := helpers.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error:": err.Error()})
			return
		}
		userID, err := strconv.ParseUint(claims.Id, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
			return
		}

		user, err := u.service.GetUserByID(uint(userID))
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		} else if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if user.UserType != "ADMIN" {
			c.JSON(http.StatusBadRequest, gin.H{"error:": "Invalid User"})
		}
		c.JSON(http.StatusOK, "Ok")

	}
}
