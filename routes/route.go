package routes

import (
	wargacontrollers "desaku-api/controllers/warga"

	"github.com/gin-gonic/gin"
)

func Route() {
	r := gin.Default()

	// --WARGA--
	// CREDENTIALS
	r.POST("/warga/login", wargacontrollers.LoginWarga)
	r.POST("/warga/register", wargacontrollers.RegisterWarga)
	r.GET("/warga/profile/:id", wargacontrollers.ProfileWargaOne) //ID WARGA

	// WIFI
	r.GET("/warga/wifi/status/:id", wargacontrollers.WifiWargaStatus) //ID WARGA

	r.Run()
}