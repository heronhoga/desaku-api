package warga

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"desaku-api/databases"
)

func GetIuranForThisMonth(c *gin.Context) {
	// Get the current year and month
	id := c.Param("id")
	now := time.Now()
	year, month, _ := now.Date()

	var iuranThisMonth []struct {
		IdIuran string `gorm:"column:id_iuran"`
		IdWarga string `gorm:"column:id_warga"`
		JumlahIuran string `gorm:"column:jumlah_iuran"`
		TanggalIuran string `gorm:"column:tanggal_iuran"`
	}

	result := databases.DB.Raw("SELECT id_iuran, id_warga, jumlah_iuran, tanggal_iuran FROM iuran WHERE id_warga = ? AND MONTH(tanggal_iuran) = ? AND YEAR(tanggal_iuran) = ?", id, month, year).Scan(&iuranThisMonth)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if len(iuranThisMonth) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ada iuran tersedia"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "data": iuranThisMonth})
}


