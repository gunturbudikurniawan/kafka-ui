package controller

import (
	"go-clean-architecture/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

// userController mengelola endpoint untuk pengguna
type userController struct {
	userService user.Service
}

// NewUserController membuat instance userController baru
func NewUserController(userService user.Service) *userController {
	return &userController{userService}
}

// RegisterUsersInput godoc
// @Summary Mendaftarkan banyak pengguna baru
// @Description Endpoint untuk registrasi batch pengguna
// @Tags Users
// @Accept json
// @Produce json
// @Param input body user.RegisterUsersInput true "Input untuk registrasi batch pengguna"
// @Success 200 {object} []user.User
// @Failure 400 {object} map[string]string
// @Router /users [post]
func (h *userController) RegisterUsersInput(c *gin.Context) {
	var input user.RegisterUsersInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, err := h.userService.RegisterUsersInput(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": users})
}

// RegisterUserInput godoc
// @Summary Mendaftarkan pengguna baru
// @Description Endpoint untuk registrasi pengguna
// @Tags Users
// @Accept json
// @Produce json
// @Param input body user.RegisterUserInput true "Input untuk registrasi pengguna"
// @Success 200 {object} user.User
// @Failure 400 {object} map[string]string
// @Router /user [post]
func (h *userController) RegisterUserInput(c *gin.Context) {
	var input user.RegisterUserInput
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user, err := h.userService.RegisterUserInput(input)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
