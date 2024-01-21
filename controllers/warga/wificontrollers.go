package warga

import (
	"github.com/gin-gonic/gin"
	"desaku-api/databases"
	"desaku-api/models"
	"net/http"
)

//CHECK WIFI SUBSCRIPTION STATUS
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

//DAFTAR WIFI
func DaftarWifiWarga(c *gin.Context) {
	var wifiData struct {
		IdWarga string `json:"id_warga"`
	}

	if err := c.ShouldBindJSON(&wifiData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if wifiData.IdWarga == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id_warga is required"})
		return
	}

	var existingWarga string
	querySearchWarga := databases.DB.Raw("SELECT id_warga FROM warga WHERE id_warga = ?", wifiData.IdWarga).Scan(&existingWarga)

	if querySearchWarga.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Warga tidak valid"})
		return
	} else if querySearchWarga.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": querySearchWarga.Error.Error()})
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

//UNSUBSCRIBE WIFI
func PutusWifiWarga(c *gin.Context) {
	var wifiData struct {
		IdWarga string `json:"id_warga"`
	}

	if err := c.ShouldBindJSON(&wifiData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if wifiData.IdWarga == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id_warga is required"})
		return
	}

	var existingWarga string
	querySearchWarga := databases.DB.Raw("SELECT id_warga FROM daftar_pelanggan_wifi WHERE id_warga = ?", wifiData.IdWarga).Scan(&existingWarga)

	if querySearchWarga.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Warga tidak valid"})
		return
	} else if querySearchWarga.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": querySearchWarga.Error.Error()})
		return
	}

	sql := "UPDATE daftar_pelanggan_wifi SET status = ? WHERE id_warga = ?"
	result := databases.DB.Exec(sql, "prosesputus", wifiData.IdWarga)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"data": map[string]string{
			"id_warga": wifiData.IdWarga,
			"status": "Sedang diproses",
		},
	})
}
