package admin

import (
	"desaku-api/databases"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllPajak (c *gin.Context) {
	year := c.Param("year")

	var WargaPajak []struct {
		IdPajak string `json:"id_pajak"`
		Tahun string `json:"tahun"`
		TotalTagihanPajak string `json:"total_tagihan_pajak"`
		StatusBayar string `json:"status_bayar"`
		Nama string `json:"nama"`
		Alamat string `json:"alamat"`
	}

	tagihanPajakQuery := databases.DB.Raw(`SELECT pajak.id_pajak,
	pajak.tahun, 
	pajak.total_tagihan_pajak, pajak.status_bayar, warga.nama, 
	warga.alamat FROM pajak INNER JOIN warga ON pajak.id_warga = 
	warga.id_warga WHERE YEAR(pajak.tahun) = ?`, year).Scan(&WargaPajak)

	if tagihanPajakQuery.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tagihanPajakQuery.Error.Error()})
		return
	}

	if tagihanPajakQuery.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No record found"})
	}

	c.JSON(http.StatusOK, gin.H{"data": WargaPajak})
}

func GetSpecificPajak(c *gin.Context) {
	id := c.Param("id") //ID PAJAK

	var WargaPajak struct {
		IdPajak string `json:"id_pajak"`
		Tahun string `json:"tahun"`
		TotalTagihanPajak string `json:"total_tagihan_pajak"`
		StatusBayar string `json:"status_bayar"`
		Nama string `json:"nama"`
		Alamat string `json:"alamat"`
	}

	tagihanPajakQuery := databases.DB.Raw(`SELECT pajak.id_pajak,
	pajak.tahun, 
	pajak.total_tagihan_pajak, pajak.status_bayar, warga.nama, 
	warga.alamat FROM pajak INNER JOIN warga ON pajak.id_warga = 
	warga.id_warga WHERE id_pajak = ?`, id).Scan(&WargaPajak)

	if tagihanPajakQuery.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": tagihanPajakQuery.Error.Error()})
		return
	}

	if tagihanPajakQuery.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No record found"})
	}

	c.JSON(http.StatusOK, gin.H{"data": WargaPajak})
}

