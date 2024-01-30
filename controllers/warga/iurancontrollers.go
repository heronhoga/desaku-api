package warga

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"desaku-api/databases"
	"strconv"
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

func BayarIuranWarga(c *gin.Context) {
	//ID IURAN
	id := c.Param("id")

	//VARIABLE FOR QUERY
	var iuranStruct struct {
		IdIuran string `gorm:"column:id_iuran"`
		IdWarga string `gorm:"column:id_warga"`
		JumlahIuran string `gorm:"column:jumlah_iuran"`
		TanggalIuran string `gorm:"column:tanggal_iuran"`
		Saldo string `gorm:"column:saldo"`
	}

	//QUERY TAGIHAN
	result := databases.DB.Raw(`
	SELECT iuran.id_iuran, iuran.id_warga, iuran.jumlah_iuran, iuran.tanggal_iuran, warga.saldo 
	FROM iuran 
	INNER JOIN warga ON iuran.id_warga = warga.id_warga
	WHERE iuran.id_iuran = ?
	`, id).Scan(&iuranStruct)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No record found"})
		return
	}

	//CONVERT TO INTEGER - jumlah_iuran
	jumlahIuranInt, err := strconv.Atoi(iuranStruct.JumlahIuran)
	if err != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"error": "Error converting jumlahIuran to int"})
	    return
	}

	//CONVERT TO INTEGER - saldo
	saldoInt, err := strconv.Atoi(iuranStruct.Saldo)
	if err != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"error": "Error converting saldo to int"})
	    return
	}

	//CONDITIONS FOR PAYMENT
	if jumlahIuranInt > saldoInt {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Saldo tidak mencukupi",
											"statusCode": http.StatusBadRequest,
											"paymentStatus": "failed"})
		return
	}

	//UPDATE SALDO
	result = databases.DB.Exec("UPDATE warga SET saldo = saldo - ? WHERE id_warga = ?", jumlahIuranInt, iuranStruct.IdWarga)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	//UPDATE STATUS
	result = databases.DB.Exec("UPDATE iuran SET status = 'paid' WHERE id_iuran = ?", id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "paymentStatus": "success"})

}