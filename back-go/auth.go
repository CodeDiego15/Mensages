package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Genera un código de verificación aleatorio
func generateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(999999)
	return fmt.Sprintf("%06d", code)
}

func handleLogin(c *gin.Context) {
	var requestBody struct {
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Generar un código de verificación
	verificationCode := generateVerificationCode()

	if requestBody.Email != "" {
		err := sendVerificationEmail(requestBody.Email, verificationCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
			return
		}
		// Almacena el código en la base de datos para su posterior verificación
	} else if requestBody.Phone != "" {
		err := sendSMS(requestBody.Phone, verificationCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send SMS"})
			return
		}
		// Almacena el código en la base de datos para su posterior verificación
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification code sent"})
}

func handleVerify(c *gin.Context) {
	var requestBody struct {
		VerificationCode string `json:"verificationCode"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Verifica el código con el almacenado en la base de datos
	// ...

	c.JSON(http.StatusOK, gin.H{"message": "Verification successful"})
}
