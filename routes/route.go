package routes

import (
	wargacontrollers "desaku-api/controllers/warga"

	"github.com/gin-gonic/gin"
)

func Route() {
	r := gin.Default()

	// --WARGA--
	// CREDENTIALS
	r.POST("/warga/login", wargacontrollers.LoginWarga) //REQUEST BODY
	r.POST("/warga/register", wargacontrollers.RegisterWarga) //REQUEST BODY
	r.GET("/warga/profile/:id", wargacontrollers.ProfileWargaOne) //ID WARGA

	// --WIFI--
	// WIFI STATUS
	r.GET("/warga/wifi/status/:id", wargacontrollers.WifiWargaStatus) //ID WARGA
	// DAFTAR WIFI
	r.POST("/warga/wifi/daftar", wargacontrollers.DaftarWifiWarga) //REQUEST BODY
	// PUTUS WIFI
	r.PUT("/warga/wifi/putus", wargacontrollers.PutusWifiWarga) //REQUEST BODY
	// BAYAR WIFI
	r.PUT("/warga/wifi/bayar/:id", wargacontrollers.BayarWifiWarga) //ID TAGIHAN
	// DAPATKAN TAGIHAN
	r.Run()
}