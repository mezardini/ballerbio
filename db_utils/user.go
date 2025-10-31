package db_utils

import (
	"log"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"github.com/golang-jwt/jwt/v5"
	"time"
	"github.com/joho/godotenv"
	"os"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"-"`
	Email    string `gorm:"unique;not null" json:"email"`
}

type CreateUserInput struct {
	gorm.Model
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

func GetUserByID(db *gorm.DB, userID uint) *User {
	var user User
	result := db.Find(&user, userID)
	if result.Error != nil {
		return nil
	}
	return &user
}

func GetUserByEmail(db *gorm.DB, userEmail string) *User {
	var user User
	result := db.Where("email = ?", userEmail).Find(&user)
	if result.Error != nil {
		return nil
	}
	return &user
}

func CreateUser(db *gorm.DB, user *User) error {
	result := db.Create(user)
	return result.Error
}

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID uint `json:"user_id"`
}


func (h *DBHandler) LoginUserGinHandler(c *gin.Context) {

	env_err := godotenv.Load()
		if env_err != nil {
			// This is useful for production where you rely only on system environment variables
			log.Println("Note: Could not find .env file, assuming environment variables are set globally.")
		}

	var jwtKey = []byte(os.Getenv("SECRET_KEY")) 

	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := GetUserByEmail(h.DB, input.Email)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password."})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password."})
		return
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	claims := &JWTClaims{ UserID : user.ID, RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(expiresAt)},}


	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token_key, err := token.SignedString(jwtKey)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful.", "user": user, "token": token_key})
}

func (h *DBHandler) CreateUserGinHandler(c *gin.Context) {
	var input CreateUserInput

	// 1. Bind JSON data to the input struct and validate required fields
	if err := c.ShouldBindJSON(&input); err != nil {
		// Log the error for debugging
		log.Printf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required and must not be empty."})
		return
	}

	check_if_user_exists := GetUserByEmail(h.DB, input.Email)
	// if err != nil {
	// 	log.Printf("Error checking user existence: %v", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not verify user."})
	// 	return
	// }

	if check_if_user_exists.ID != 0 && check_if_user_exists != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with given Email already exists."})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Password hashing error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process password."})
		return
	}

	user := User{
		Username: input.Username,
		Password: string(hashedPassword),
		Email:    input.Email,
	}

	// 3. Call the database function
	if err := CreateUser(h.DB, &user); err != nil {
		log.Printf("Database create error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create profile."})
		return
	}

	// 4. Respond with the newly created profile (including the new ID)
	c.JSON(http.StatusCreated, user)
}

func (h *DBHandler) GetUserByIDGinHandler(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID format.",
		})
		return
	}
	user := GetUserByID(h.DB, uint(userID))

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found.",
		})
		return
	}		
	c.JSON(http.StatusOK, user)
}