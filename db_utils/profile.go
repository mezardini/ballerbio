package db_utils

import (
	"ballerbio/utils"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Profile struct {
	gorm.Model
	FirstName    string        `json:"first_name"`
	LastName     string        `json:"last_name"`
	Dob          time.Time     `json:"dob"`
	Position     string        `json:"position"`
	Height       float64       `json:"height"`
	Weight       float64       `json:"weight"`
	Bio          string        `json:"bio"`
	Location     string        `json:"location"`
	Nationality  string        `json:"nationality"`
	Slug         string        `gorm:"not null" json:"slug"`
	UserID       uint          `gorm:"unique;not null" json:"user_id"`
	User         User          `json:"user"`
	Skills       []Skill       `json:"skills"`
	Achievements []Achievement `json:"achievements"`
	Injuries     []Injury      `json:"injuries"`
	SocialLinks  []SocialLink  `json:"social_links"`
	ClubProfiles []ClubProfile `json:"club_profiles"`
	SeasonStats  []SeasonStat  `json:"season_stats"`
}

type CreateProfileInput struct {
	gorm.Model

	FirstName   string    `json:"first_name" binding:"required"`
	LastName    string    `json:"last_name" binding:"required"`
	Dob         time.Time `json:"dob" binding:"required"`
	Position    string    `json:"position" binding:"required"`
	Height      float64   `json:"height" binding:"required"`
	Weight      float64   `json:"weight" binding:"required"`
	Bio         string    `json:"bio" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	Nationality string    `json:"nationality" binding:"required"`
	UserID      uint      `json:"user_id" binding:"required"`
}

func GetProfileByID(db *gorm.DB, profileID uint) (Profile, error) {
	var profile Profile

	// Preload ALL relationships
	result := db.First(&profile, profileID)

	return profile, result.Error
}

func GetProfiles(db *gorm.DB) ([]Profile, error) {
	var profiles []Profile

	// Preload ALL relationships
	result := db.
		Preload("User").
		Preload("Skills").
		Preload("Achievements").
		Preload("Injuries").
		Preload("SocialLinks").
		Preload("ClubProfiles").
		Preload("SeasonStats").
		Find(&profiles)

	return profiles, result.Error
}

func (h *DBHandler) GetProfilesGinHandler(c *gin.Context) {
	// 1. Call the database function (no package prefix needed for GetProfiles)
	profiles, err := GetProfiles(h.DB)

	if err != nil {
		log.Printf("Error fetching profiles: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not retrieve profiles.",
		})
		return
	}

	// 2. Respond with the fetched data
	c.JSON(http.StatusOK, profiles)
}

func GetProfile(db *gorm.DB, profileID uint, slug string) (Profile, error) {
	var profile Profile

	// Preload ALL relationships
	result := db.
		Preload("User").
		Preload("Skills").
		Preload("Achievements").
		Preload("Injuries").
		Preload("SocialLinks").
		Preload("ClubProfiles").
		Preload("SeasonStats").
		Where("id = ? AND slug = ?", profileID, slug).
		First(&profile)

	// utils.SendEmail(profile.User.Email, "Profile Viewed", "Your profile was just viewed.")

	return profile, result.Error
}

func (h *DBHandler) GetProfileByIDGinHandler(c *gin.Context) {
	// 1. Extract and validate ID from URL parameter
	idParam := c.Param("id")
	slugParam := c.Param("slug")

	profileID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid profile ID format.",
		})
		return
	}

	// 2. Call the database function
	profile, err := GetProfile(h.DB, uint(profileID), string(slugParam))

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Profile not found.",
			})
			return
		}

		log.Printf("Error fetching profile ID %d: %v", profileID, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not retrieve profile.",
		})
		return
	}

	// 3. Respond with the fetched data
	c.JSON(http.StatusOK, profile)
}

func CreateProfile(db *gorm.DB, profile *Profile) error {
	result := db.Create(profile)
	return result.Error
}

func (h *DBHandler) CreateProfileGinHandler(c *gin.Context) {
	var input CreateProfileInput

	// 1. Bind JSON data to the input struct and validate required fields
	if err := c.ShouldBindJSON(&input); err != nil {
		// Log the error for debugging
		log.Printf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.UserID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UserID is required and must be non-zero."})
		return
	}

	check_if_user_exists := GetUserByID(h.DB, input.UserID)
	// if err != nil {
	// 	log.Printf("Error checking user existence: %v", err)
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not verify user."})
	// 	return
	// }

	if check_if_user_exists.ID == 0 || check_if_user_exists == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with given UserID does not exist."})
		return
	}

	create_slug := utils.ProfileSlugify(input.FirstName, input.LastName)

	profile := Profile{
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		Dob:         input.Dob,
		Position:    input.Position,
		Height:      input.Height,
		Weight:      input.Weight,
		Bio:         input.Bio,
		Location:    input.Location,
		Nationality: input.Nationality,
		Slug:        create_slug,
		UserID:      input.UserID,
	}

	// 3. Call the database function
	if err := CreateProfile(h.DB, &profile); err != nil {
		log.Printf("Database create error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create profile."})
		return
	}

	// 4. Respond with the newly created profile (including the new ID)
	c.JSON(http.StatusCreated, profile)
}
