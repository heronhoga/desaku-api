package warga

import (
	"github.com/gin-gonic/gin"
	"desaku-api/databases"
	"desaku-api/models"
	"net/http"
)
func WifiWargaStatus(c *gin.Context) {
	id := c.Param("id")

	var langgananWifi models.LanggananWifi

	result := databases.DB.Raw("SELECT * FROM daftar_pelanggan_wifi WHERE id_warga = ?", id).Scan(&langgananWifi)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No record found with the provided ID"})
		return
		} else if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"data":       langgananWifi,
	})
}

func DaftarWifiWarga(c *gin.Context) {
	var wifiData struct {
		IdWarga string `json:"id_warga"`
	}

	if err := c.ShouldBindJSON(&wifiData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var langgananWifi models.LanggananWifi
	result := databases.DB.Raw("SELECT * FROM daftar_pelanggan_wifi WHERE id_warga = ?", wifiData.IdWarga).Scan(&langgananWifi)

	if result.RowsAffected > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "User already subscribed"})
		return
	}

	sql := "INSERT INTO daftar_pelanggan_wifi (id_warga, status) VALUES (?, ?)"
	result = databases.DB.Exec(sql, wifiData.IdWarga, "proses")
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"statusCode": http.StatusCreated,
		"data": map[string]string{
			"id_warga": wifiData.IdWarga,
			"status": "Sedang diproses",
		},
	})
}
