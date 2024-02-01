package admin

import (
	"desaku-api/databases"
	"net/http"
	"time"
	"fmt"
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

func CreateIuran (c *gin.Context) {
	month := c.Param("month")
	
	var dataIuran[] struct {
		IdWarga string `json:"id_warga"`
	}

	wargaQuery := databases.DB.Raw(
		`SELECT id_warga FROM warga`).Scan(&dataIuran)
	
	if wargaQuery.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": wargaQuery.Error.Error()})
		return
	}

	if wargaQuery.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No record found"})
		return
	}

	currentYear := time.Now().Year()
	date := fmt.Sprintf("%d-%s-01", currentYear, month)

	for _, warga := range dataIuran {
		insertIuran := databases.DB.Exec(`INSERT INTO iuran (id_warga, jumlah_iuran, tanggal_iuran) 
		VALUES (?, ?, ?)`, warga.IdWarga, 100000, date)

		if insertIuran.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": insertIuran.Error.Error()})
			return
		}

	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Data created successfully"})
}

func DeleteIuran (c *gin.Context) {
	id := c.Param("id")

	deleteQuery := databases.DB.Exec("DELETE FROM iuran WHERE id_iuran = ?", id)
	if deleteQuery.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": deleteQuery.Error.Error()})
		return
	}

	if deleteQuery.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No record found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data deleted successfully"})
}