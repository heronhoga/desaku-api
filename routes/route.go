package routes

import (
	profilecontrollers "desaku-api/controllers/warga"

	"github.com/gin-gonic/gin"
)

func Route() {
	r := gin.Default()

	// --WARGA--
	// CREDENTIALS
	r.POST("/warga/login", profilecontrollers.LoginWarga)
	r.POST("/warga/register", profilecontrollers.RegisterWarga)


	r.Run()
}