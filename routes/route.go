package routes

import (
	wargacontrollers "desaku-api/controllers/warga"
	admincontrollers "desaku-api/controllers/admin"
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
	r.GET("/warga/wifi/tagihan/:id", wargacontrollers.GetTagihanWarga) //ID WARGA

	// --PAJAK--
	// DAPATKAN TAGIHAN PAJAK
	r.GET("/warga/pajak/tagihan/:id", wargacontrollers.GetTagihanPajakWarga) //ID WARGA
	// BAYAR PAJAK
	r.PUT("/warga/pajak/bayar/:id", wargacontrollers.BayarPajakWarga) //ID TAGIHAN

	// --ARTIKEL--
	// DAPATKAN SEMUA ARTIKEL
	r.GET("/warga/artikel", wargacontrollers.GetAllArtikel)

	// --E-PASAR--
	// DAPATKAN SEMUA E-PASAR
	r.GET("/warga/epasar", wargacontrollers.GetAllEpasar)

	//--LISTRIK--
	// DAPATKAN TAGIHAN LISTRIK
	r.GET("/warga/listrik/tagihan/:id", wargacontrollers.GetTagihanListrikWarga) //ID WARGA
	// BAYAR LISTRIK
	r.PUT("/warga/listrik/bayar/:id", wargacontrollers.BayarListrikWarga) //ID TAGIHAN

	//--IURAN--
	// DAPATKAN TAGIHAN IURAN
	r.GET("/warga/iuran/tagihan/:id", wargacontrollers.GetIuranForThisMonth) //ID WARGA
	// BAYAR IURAN
	r.PUT("/warga/iuran/bayar/:id", wargacontrollers.BayarIuranWarga) //ID IURAN


	//--ADMIN--
	// AUTH
	r.POST("/admin/login", admincontrollers.LoginAdmin) //REQUEST BODY
	r.Run()
}