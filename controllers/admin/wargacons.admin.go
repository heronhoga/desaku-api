package admin

import (
	"net/http"
	"desaku-api/databases"
	"github.com/gin-gonic/gin"
)

func GetAllWarga (c *gin.Context) {
	var warga[] struct {
		IdWarga string `json:"id_warga"`
		Nama string `json:"nama"`
		TanggalLahir string `json:"tanggal_lahir"`
		JenisKelamin string `json:"jenis_kelamin"`
		Nik string `json:"nik"`
		Alamat string `json:"alamat"`
		Saldo string `json:"saldo"`
	}

	wargaQuery := databases.DB.Raw(`SELECT id_warga, nama, tanggal_lahir, jenis_kelamin, nik, alamat, saldo FROM warga`).Scan(&warga)
	if wargaQuery.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": wargaQuery.Error.Error()})
		return
	}

if wargaQuery.RowsAffected == 0 {
	c.JSON(http.StatusNotFound, gin.H{"error": "No record found"})
}

	c.JSON(http.StatusOK, gin.H{"data": warga})
}