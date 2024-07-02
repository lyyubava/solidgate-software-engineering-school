package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/lyyubava/solidgate-software-engineering-school.git/models"
	"net/http"
)

type EmailInput struct {
	Email string `json:"email" binding:"required"`
}

func Subscribe(c *gin.Context) {
	var emailInput EmailInput
	if err := c.ShouldBindJSON(&emailInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	emailDB := models.Email{}
	models.DB.Model(models.Email{}).Where("email = ?", emailInput.Email).First(&emailDB)
	if emailDB.Email == emailInput.Email {
		c.JSON(http.StatusConflict, gin.H{"error": "email already subscribed"})
		return
	}

	email := models.Email{Email: emailInput.Email}
	models.DB.Create(&email)
	c.JSON(http.StatusOK, gin.H{"message": "email was successfully subscribed"})
}
