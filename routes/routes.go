package routes

import (
	"ballerbio/db_utils"
	"ballerbio/middleware"
	"log"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes() {

	db, err := db_utils.ConnectAndMigrate()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Inject the DB connection into the handler struct
	handler := &db_utils.DBHandler{DB: db} 

	router := gin.Default()

	authorized := router.Group("/api")
	authorized.Use(auth.AuthMiddleware())
	{
		authorized.POST("/profiles/create", handler.CreateProfileGinHandler)
		authorized.POST("/skills/add", handler.AddSkillToProfileGinHandler)
		authorized.POST("/achievements/add", handler.AddAchievementToProfileGinHandler)
		authorized.POST("/injury/add", handler.AddInjuryToProfileGinHandler)
		authorized.POST("/sociallink/add", handler.AddSocialLinkToProfileGinHandler)
		authorized.POST("/clubprofile/add", handler.AddClubProfileToProfileGinHandler)
		authorized.POST("/seasonstats/add", handler.AddSeasonStatToProfileGinHandler)
	}
	

	router.GET("/profiles", handler.GetProfilesGinHandler)
	router.GET("/profiles/:id/:slug", handler.GetProfileByIDGinHandler)
	router.GET("/skills/:id", handler.GetPlayerSkillsGinHandler)
	// router.POST("/profiles/create", handler.CreateProfileGinHandler)
	router.GET("/users/:id", handler.GetUserByIDGinHandler)
	router.POST("/users/create", handler.CreateUserGinHandler)
	router.POST("/users/login", handler.LoginUserGinHandler)
	// router.POST("/skills/add", handler.AddSkillToProfileGinHandler)

	router.Run("localhost:8081")
}