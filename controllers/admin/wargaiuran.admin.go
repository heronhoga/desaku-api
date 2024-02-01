package admin

import (
	"desaku-api/databases"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)
func GetAllIuran (c *gin.Context) {
	
	month := c.Param("month")
	year := time.Now().Year()

	var WargaIuran []struct {
		IdIuran string `json:"id_iuran"`
		JumlahIuran string `json:"jumlah_iuran"`
		TanggalIuran string `json:"tanggal_iuran"`
		Status string `json:"status"`
		Nama string `json:"nama"`
		Alamat string `json:"alamat"`
	}

	iuranQuery := databases.DB.Raw(
		`SELECT iuran.id_iuran, iuran.jumlah_iuran, iuran.tanggal_iuran,
		iuran.status, warga.nama, warga.alamat FROM iuran INNER JOIN warga
		ON iuran.id_warga = warga.id_warga WHERE MONTH(iuran.tanggal_iuran) = ?
		AND YEAR(iuran.tanggal_iuran) = ?`, month, year).Scan(&WargaIuran)
	
	if iuranQuery.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": iuranQuery.Error.Error()})
		return
	}

	if iuranQuery.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No record found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": WargaIuran})
}