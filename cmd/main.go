// @title           VatanSoft Hospital API
// @version         1.0
// @description     Backend case project for VatanSoft.

// @contact.name   Efecan Yılmazdemir
// @contact.email  efecan.yilmazdemir@gmail.com

// @host      localhost:8080
// @BasePath  /

package main

import (
	"net/http"

	"github.com/efecan/vatansoft-case/config"
	"github.com/efecan/vatansoft-case/controllers"
	_ "github.com/efecan/vatansoft-case/docs"
	"github.com/efecan/vatansoft-case/middlewares"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()
	config.InitRedis()

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.POST("/auth/request-password-reset", controllers.RequestPasswordReset)
	r.POST("/auth/reset-password", controllers.ResetPassword)
	r.POST("/hospitals/register", middlewares.RequireAuth, middlewares.RequireAdmin, controllers.HospitalRegister)
	r.GET("/hospitals", middlewares.RequireAuth, controllers.GetHospitals)
	r.POST("/users", middlewares.RequireAuth, middlewares.RequireAdmin, controllers.CreateUser)
	r.GET("/users", middlewares.RequireAuth, controllers.GetUsers)
	r.GET("/listusers", middlewares.RequireAuth, controllers.ListUsers)
	r.PUT("/users/:id", middlewares.RequireAuth, middlewares.RequireAdmin, controllers.UpdateUser)
	r.DELETE("/users/:id", middlewares.RequireAuth, middlewares.RequireAdmin, controllers.DeleteUser)
	r.POST("/departments", middlewares.RequireAuth, middlewares.RequireAdmin, controllers.CreateDepartment)
	r.GET("/departments", middlewares.RequireAuth, controllers.GetDepartments)
	r.GET("/departments/:id/doctors", middlewares.RequireAuth, controllers.GetDoctorsByDepartment)
	r.GET("/cities", middlewares.RequireAuth, controllers.GetCities)
	r.GET("/profession-groups", controllers.GetProfessionGroups)

	r.GET("/me", middlewares.RequireAuth, func(c *gin.Context) {
		userID := c.GetInt("userID")
		userName := c.GetString("userName")

		c.JSON(http.StatusOK, gin.H{
			"userID":   userID,
			"userName": userName,
		})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := config.GetEnv("PORT", "8080")
	r.Run(":" + port)
}
