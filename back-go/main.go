package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

var (
	db        *gorm.DB
	jwtSecret = []byte("your_secret_key")
)

type User struct {
	ID               uint   `gorm:"primary_key"`
	Email            string `gorm:"unique_index"`
	Phone            string `gorm:"unique_index"`
	VerificationCode string
	IsVerified       bool
}

type Message struct {
	ID        uint `gorm:"primary_key"`
	UserID    uint
	Content   string
	CreatedAt time.Time
}

func main() {
	var err error
	db, err = gorm.Open("sqlite3", "messaging.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	db.AutoMigrate(&User{}, &Message{})

	r := gin.Default()
	r.POST("/login", login)
	r.POST("/verify", verify)
	r.POST("/send", sendMessage)
	r.GET("/messages", getMessages)

	r.Run(":8080")
}

func login(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var existingUser User
	if user.Email != "" {
		if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}
	} else if user.Phone != "" {
		if err := db.Where("phone = ?", user.Phone).First(&existingUser).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or phone required"})
		return
	}

	// Generate verification code
	verificationCode := generateVerificationCode()
	existingUser.VerificationCode = verificationCode
	db.Save(&existingUser)

	if user.Email != "" {
		err := sendVerificationEmail(user.Email, verificationCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send verification email"})
			return
		}
	} else if user.Phone != "" {
		err := sendSMS(user.Phone, verificationCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send SMS"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Verification code sent"})
}

func verify(c *gin.Context) {
	var requestBody struct {
		Email            string `json:"email"`
		Phone            string `json:"phone"`
		VerificationCode string `json:"verificationCode"`
	}
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var user User
	if requestBody.Email != "" {
		if err := db.Where("email = ? AND verification_code = ?", requestBody.Email, requestBody.VerificationCode).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid verification code"})
			return
		}
	} else if requestBody.Phone != "" {
		if err := db.Where("phone = ? AND verification_code = ?", requestBody.Phone, requestBody.VerificationCode).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid verification code"})
			return
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or phone required"})
		return
	}

	user.IsVerified = true
	user.VerificationCode = ""
	db.Save(&user)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, _ := token.SignedString(jwtSecret)

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func sendMessage(c *gin.Context) {
	var tokenString string
	tokenString = c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
		return
	}

	tokenString = tokenString[len("Bearer "):]

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	userID := uint(claims["user_id"].(float64))
	var msg Message
	if err := c.BindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	msg.UserID = userID
	msg.CreatedAt = time.Now()
	db.Create(&msg)

	c.JSON(http.StatusOK, gin.H{"message": "Message sent"})
}

func getMessages(c *gin.Context) {
	var messages []Message
	db.Order("created_at desc").Find(&messages)
	c.JSON(http.StatusOK, messages)
}

func generateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(999999)
	return fmt.Sprintf("%06d", code)
}
