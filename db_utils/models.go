package db_utils

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

// --- Structs (Unchanged) ---
type DBHandler struct {
	DB *gorm.DB
}

type Skill struct {
	gorm.Model
	Name      string  `json:"skill_name"`
	Level     string  `json:"level"`
	ProfileID uint    `gorm:"not null" json:"profile_id"`
	Profile   Profile `json:"profile"`
}

type AddSkill struct {
	gorm.Model

	Name      string  `json:"skill_name" binding:"required"`
	Level     string  `json:"level" binding:"required"`
	ProfileID uint    `json:"profile_id" binding:"required"`
	Profile   Profile `json:"profile" binding:"required"`
}

type Achievement struct {
	gorm.Model

	Title        string     `gorm:"size:100" json:"title"`
	Description  string     `json:"description"`
	DateAchieved *time.Time `json:"date_achieved"`
	ProfileID    uint       `gorm:"not null" json:"profile_id"`
	Profile      Profile    `json:"profile"`
}

type AddAchievement struct {
	gorm.Model	
	Title        string     `gorm:"size:100" json:"title" binding:"required"`
	Description  string     `json:"description" binding:"required"`
	DateAchieved *time.Time `json:"date_achieved" binding:"required"`
	ProfileID    uint       `gorm:"not null" json:"profile_id" binding:"required"`
	Profile      Profile    `json:"profile" binding:"required"`
}

type Injury struct {
	gorm.Model

	InjuryType  string     `gorm:"size:100" json:"injury_type"`
	Description string     `json:"description"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	ProfileID   uint       `gorm:"not null" json:"profile_id"`
	Profile     Profile    `json:"profile"`
}

type AddInjury struct {
	gorm.Model

	InjuryType  string     `gorm:"size:100" json:"injury_type" binding:"required"`
	Description string     `json:"description" binding:"required"`
	StartDate   *time.Time `json:"start_date" binding:"required"`
	EndDate     *time.Time `json:"end_date"`
	ProfileID   uint       `gorm:"not null" json:"profile_id" binding:"required"`
	Profile     Profile    `json:"profile" binding:"required"`
}

type SocialLink struct {
	gorm.Model

	Platform  string  `gorm:"size:50" json:"platform"`
	URL       string  `gorm:"size:200" json:"url"`
	ProfileID uint    `gorm:"not null" json:"profile_id"`
	Profile   Profile `json:"profile"`
}

type AddSocialLink struct {
	gorm.Model

	Platform  string  `gorm:"size:50" json:"platform" binding:"required"`
	URL       string  `gorm:"size:200" json:"url" binding:"required"`
	ProfileID uint    `gorm:"not null" json:"profile_id" binding:"required"`
	Profile   Profile `json:"profile" binding:"required"`
}

const (
	ContractTypePermanent = "Permanent"
	ContractTypeLoan      = "Loan"
	ContractTypeTrial     = "Trial"
)

type ClubProfile struct {
	gorm.Model

	ProfileID       *uint      `json:"profile_id"`
	ClubName        string     `gorm:"size:100" json:"club_name"`
	ClubLeague      string     `gorm:"size:100" json:"club_league"`
	ClubCountry     string     `gorm:"size:100" json:"club_country"`
	StartYear       *time.Time `json:"start_year"`
	EndYear         *time.Time `json:"end_year"`
	IsPresentClub   bool       `gorm:"default:false" json:"is_present_club"`
	ClubAppearances *int32     `json:"club_appearances"`
	ClubGoals       *int32     `json:"club_goals"`
	ClubAssists     *int32     `json:"club_assists"`
	ContractType    string     `gorm:"size:20;default:'Permanent'" json:"contract_type"`
	Profile         Profile    `json:"profile"`
}

type AddClubProfile struct {
	gorm.Model

	ProfileID       *uint      `json:"profile_id" binding:"required"`
	ClubName        string     `gorm:"size:100" json:"club_name" binding:"required"`
	ClubLeague      string     `gorm:"size:100" json:"club_league" binding:"required"`	
	ClubCountry     string     `gorm:"size:100" json:"club_country" binding:"required"`
	StartYear       *time.Time `json:"start_year" binding:"required"`
	EndYear         *time.Time `json:"end_year"`
	IsPresentClub   bool       `gorm:"default:false" json:"is_present_club"`
	ClubAppearances *int32     `json:"club_appearances"`
	ClubGoals       *int32     `json:"club_goals"`
	ClubAssists     *int32     `json:"club_assists"`
	ContractType    string     `gorm:"size:20;default:'Permanent'" json:"contract_type" binding:"required"`
	Profile         Profile    `json:"profile" binding:"required"`
}

type SeasonStat struct {
	gorm.Model

	ProfileID     uint    `gorm:"not null" json:"profile_id"`
	Season        string  `gorm:"size:20" json:"season"`
	ClubName      string  `gorm:"size:100" json:"club_name"`
	LeagueName    string  `gorm:"size:100" json:"league_name"`
	Appearances   *int32  `json:"appearances"`
	Goals         *int32  `json:"goals"`
	Assists       *int32  `json:"assists"`
	MinutesPlayed *int32  `json:"minutes_played"`
	YellowCards   *int32  `json:"yellow_cards"`
	RedCards      *int32  `json:"red_cards"`
	Profile       Profile `json:"profile"`
}

type AddSeasonStat struct {
	gorm.Model

	ProfileID     uint    `gorm:"not null" json:"profile_id" binding:"required"`
	Season        string  `gorm:"size:20" json:"season" binding:"required"`
	ClubName      string  `gorm:"size:100" json:"club_name" binding:"required"`	
	LeagueName    string  `gorm:"size:100" json:"league_name" binding:"required"`
	Appearances   *int32  `json:"appearances"`
	Goals         *int32  `json:"goals"`
	Assists       *int32  `json:"assists"`
	MinutesPlayed *int32  `json:"minutes_played"`
	YellowCards   *int32  `json:"yellow_cards"`
	RedCards      *int32  `json:"red_cards"`
	Profile       Profile `json:"profile" binding:"required"`
}

func GetPlayerSkills(db *gorm.DB, profileID uint) ([]Skill, error) {
	var skills []Skill
	result := db.Preload("Profile").Find(&skills, profileID)
	return skills, result.Error
}

func (h *DBHandler) GetPlayerSkillsGinHandler(c *gin.Context) {
	idParam := c.Param("id")
	profileID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid profile ID format.",
		})
		return
	}

	profile, err := GetPlayerSkills(h.DB, uint(profileID))

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

	c.JSON(http.StatusOK, profile)
}

func AddSkillToProfile(db *gorm.DB, skill *Skill) error {
	result := db.Create(skill)
	return result.Error
}

func (h *DBHandler) AddSkillToProfileGinHandler(c *gin.Context) {
	var input AddSkill

	// 1. Bind JSON data to the input struct and validate required fields
	if err := c.ShouldBindJSON(&input); err != nil {
		// Log the error for debugging
		log.Printf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.ProfileID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ProfileID is required and must be non-zero."})
		return
	}

	check_if_profile_exists, err := GetProfileByID(h.DB, input.ProfileID)
	if err != nil {
		log.Printf("Error checking profile existence: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not verify profile."})
		return
	}
	if check_if_profile_exists.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Profile with given ID does not exist."})
		return
	}

	skill := Skill{
		Name:      input.Name,
		Level:     input.Level,
		ProfileID: input.ProfileID,
		Profile:   check_if_profile_exists,
	}

	if err := AddSkillToProfile(h.DB, &skill); err != nil {
		log.Printf("Database create error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create profile."})
		return
	}

	// 4. Respond with the newly created profile (including the new ID)
	c.JSON(http.StatusCreated, skill)
}

func AddAchievementToProfile(db *gorm.DB, achievement *Achievement) error {
	result := db.Create(achievement)
	return result.Error
}

func (h *DBHandler) AddAchievementToProfileGinHandler(c *gin.Context) {
	var input AddAchievement	

	// 1. Bind JSON data to the input struct and validate required fields
	if err := c.ShouldBindJSON(&input); err != nil {
		// Log the error for debugging
		log.Printf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.ProfileID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ProfileID is required and must be non-zero."})
		return
	}
	check_if_profile_exists, err := GetProfileByID(h.DB, input.ProfileID)
	if err != nil {
		log.Printf("Error checking profile existence: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not verify profile."})
		return
	}	
	if check_if_profile_exists.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Profile with given ID does not exist."})
		return
	}
	achievement := Achievement{
		Title:        input.Title,
		Description:  input.Description,	
		DateAchieved: input.DateAchieved,
		ProfileID:    input.ProfileID,
		Profile:      check_if_profile_exists,
	}	
	if err := AddAchievementToProfile(h.DB, &achievement); err != nil {
		log.Printf("Database create error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create achievement."})
		return
	}
	// 4. Respond with the newly created achievement (including the new ID)
	c.JSON(http.StatusCreated, achievement)
}

func AddInjuryToProfile(db *gorm.DB, injury *Injury) error {
	result := db.Create(injury)
	return result.Error
}

func (h *DBHandler) AddInjuryToProfileGinHandler(c *gin.Context) {
	var input AddInjury	
	// 1. Bind JSON data to the input struct and validate required fields
	if err := c.ShouldBindJSON(&input); err != nil {
		// Log the error for debugging
		log.Printf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.ProfileID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ProfileID is required and must be non-zero."})
		return
	}
	check_if_profile_exists, err := GetProfileByID(h.DB, input.ProfileID)
	if err != nil {
		log.Printf("Error checking profile existence: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not verify profile."})
		return
	}
	if check_if_profile_exists.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Profile with given ID does not exist."})
		return
	}
	injury := Injury{
		InjuryType:  input.InjuryType,
		Description: input.Description,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
		ProfileID:   input.ProfileID,
		Profile:     check_if_profile_exists,
	}
	if err := AddInjuryToProfile(h.DB, &injury); err != nil {
		log.Printf("Database create error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create injury."})
		return
	}
	// 4. Respond with the newly created injury (including the new ID)
	c.JSON(http.StatusCreated, injury)
}

func AddSocialLinkToProfile(db *gorm.DB, socialLink *SocialLink) error {
	result := db.Create(socialLink)
	return result.Error
}

func (h *DBHandler) AddSocialLinkToProfileGinHandler(c *gin.Context) {
	var input AddSocialLink
	// 1. Bind JSON data to the input struct and validate required fields
	if err := c.ShouldBindJSON(&input); err != nil {
		// Log the error for debugging
		log.Printf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}	
	if input.ProfileID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ProfileID is required and must be non-zero."})
		return
	}
	check_if_profile_exists, err := GetProfileByID(h.DB, input.ProfileID)
	if err != nil {
		log.Printf("Error checking profile existence: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not verify profile."})
		return
	}
	if check_if_profile_exists.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Profile with given ID does not exist."})
		return
	}
	socialLink := SocialLink{
		Platform:  input.Platform,
		URL:       input.URL,
		ProfileID: input.ProfileID,
		Profile:   check_if_profile_exists,
	}
	if err := AddSocialLinkToProfile(h.DB, &socialLink); err != nil {
		log.Printf("Database create error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create social link."})
		return
	}
	// 4. Respond with the newly created social link (including the new ID)
	c.JSON(http.StatusCreated, socialLink)
}

func AddClubProfileToProfile(db *gorm.DB, clubProfile *ClubProfile) error {
	result := db.Create(clubProfile)
	return result.Error
}

func (h *DBHandler) AddClubProfileToProfileGinHandler(c *gin.Context) {
	var input AddClubProfile
	// 1. Bind JSON data to the input struct and validate required fields
	if err := c.ShouldBindJSON(&input); err != nil {
		// Log the error for debugging
		log.Printf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.ProfileID == nil || *input.ProfileID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ProfileID is required and must be non-zero."})
		return
	}
	check_if_profile_exists, err := GetProfileByID(h.DB, *input.ProfileID)
	if err != nil {
		log.Printf("Error checking profile existence: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not verify profile."})
		return
	}
	if check_if_profile_exists.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Profile with given ID does not exist."})
		return
	}
	clubProfile := ClubProfile{
		ProfileID:       input.ProfileID,
		ClubName:        input.ClubName,	
		ClubLeague:      input.ClubLeague,
		ClubCountry:     input.ClubCountry,
		StartYear:       input.StartYear,
		EndYear:         input.EndYear,
		IsPresentClub:   input.IsPresentClub,
		ClubAppearances: input.ClubAppearances,
		ClubGoals:       input.ClubGoals,
		ClubAssists:     input.ClubAssists,
		ContractType:    input.ContractType,
		Profile:         check_if_profile_exists,
	}
	if err := AddClubProfileToProfile(h.DB, &clubProfile); err != nil {
		log.Printf("Database create error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create club profile."})
		return
	}	
	// 4. Respond with the newly created club profile (including the new ID)
	c.JSON(http.StatusCreated, clubProfile)
}

func AddSeasonStatToProfile(db *gorm.DB, seasonStat *SeasonStat) error {
	result := db.Create(seasonStat)
	return result.Error
}

func (h *DBHandler) AddSeasonStatToProfileGinHandler(c *gin.Context) {
	var input AddSeasonStat	
	// 1. Bind JSON data to the input struct and validate required fields
	if err := c.ShouldBindJSON(&input); err != nil {
		// Log the error for debugging
		log.Printf("JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if input.ProfileID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ProfileID is required and must be non-zero."})
		return
	}
	check_if_profile_exists, err := GetProfileByID(h.DB, input.ProfileID)
	if err != nil {
		log.Printf("Error checking profile existence: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not verify profile."})
		return
	}
	if check_if_profile_exists.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Profile with given ID does not exist."})
		return
	}
	seasonStat := SeasonStat{
		ProfileID:     input.ProfileID,
		Season:        input.Season,	
		ClubName:      input.ClubName,
		LeagueName:    input.LeagueName,
		Appearances:   input.Appearances,
		Goals:         input.Goals,
		Assists:       input.Assists,
		MinutesPlayed: input.MinutesPlayed,
		YellowCards:   input.YellowCards,
		RedCards:      input.RedCards,
		Profile:       check_if_profile_exists,
	}
	if err := AddSeasonStatToProfile(h.DB, &seasonStat); err != nil {
		log.Printf("Database create error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create season stat."})
		return
	}
	// 4. Respond with the newly created season stat (including the new ID)
	c.JSON(http.StatusCreated, seasonStat)
}