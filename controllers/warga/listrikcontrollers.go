package warga

import (
	"desaku-api/databases"
	"net/http"
	"strconv"
	"github.com/gin-gonic/gin"
)

func GetTagihanListrikWarga(c *gin.Context) {
	//ID WARGA
	id := c.Param("id")

	//VARIABLE FOR TAGIHAN
	var tagihanListrik []struct {
		IdListrik string `gorm:"column:id_listrik"`
		IdWarga string `gorm:"column:id_warga"`
		TotalTagihanListrik int64 `gorm:"column:total_tagihan_listrik"`
		TanggalTagihan string `gorm:"column:tanggal_tagihan"`
	}

	//GETTING ALL DATA
	result := databases.DB.Raw("SELECT id_listrik, id_warga, total_tagihan_listrik, tanggal_tagihan FROM listrik WHERE id_warga = ? AND status = 'pending'", id).Scan(&tagihanListrik)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if len(tagihanListrik) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ada tagihan listrik tersedia"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "data": tagihanListrik})
}

func BayarListrikWarga(c *gin.Context) {
	//ID TAGIHAN
	id := c.Param("id")

	//VARIABLE FOR QUERY
	var resultStruct struct {
		IdListrik string `gorm:"column:id_listrik"`
		IdWarga string `gorm:"column:id_warga"`
		TotalTagihanListrik string `gorm:"column:total_tagihan_listrik"`
		TanggalTagihan string `gorm:"column:tanggal_tagihan"`
	}

	//QUERY TAGIHAN
	result := databases.DB.Raw("SELECT id_listrik, id_warga, total_tagihan_listrik, tanggal_tagihan FROM listrik WHERE id_listrik = ?", id).Scan(&resultStruct)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tagihan not found"})
		return
	}

	//QUERY SALDO WARGA
	var saldoWargaStruct struct {
		Saldo string `gorm:"column:saldo"`
	}

	result = databases.DB.Raw("SELECT saldo FROM warga WHERE id_warga = ?", resultStruct.IdWarga).Scan(&saldoWargaStruct)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	//CONVERT TO INTEGER
	saldoInt, err := strconv.Atoi(saldoWargaStruct.Saldo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error converting saldo to int"})
		return
	}

	//CONVERT TO INTEGER
	totalTagihanListrikInt, err := strconv.Atoi(resultStruct.TotalTagihanListrik)
	if err != nil {
	    c.JSON(http.StatusInternalServerError, gin.H{"error": "Error converting totalTagihanListrik to int"})
	    return
	}

	//CONDITIONS FOR PAYMENT
	if totalTagihanListrikInt > saldoInt {
	    c.JSON(http.StatusBadRequest, gin.H{"error": "Saldo tidak mencukupi",
	                                        "statusCode": http.StatusBadRequest,
	                                        "paymentStatus": "failed"})
	    return
	}

	//UPDATE STATUS
	result = databases.DB.Exec("UPDATE listrik SET status = 'paid' WHERE id_listrik = ?", id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	//UPDATE SALDO
	result = databases.DB.Exec("UPDATE warga SET saldo = saldo - ? WHERE id_warga = ?", totalTagihanListrikInt, resultStruct.IdWarga)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "paymentStatus": "success"})

}