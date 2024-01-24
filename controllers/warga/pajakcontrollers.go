package warga

import (
	"desaku-api/databases"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetTagihanPajakWarga(c *gin.Context) {
    //ID WARGA
    id := c.Param("id")

    //VARIABLE FOR TAGIHAN
    var tagihanPajak []struct {
        IdPajak string `gorm:"column:id_pajak"`
		IdWarga string `gorm:"column:id_warga"`
		Tahun string `gorm:"column:tahun"`
		TotalTagihanPajak string `gorm:"column:total_tagihan_pajak"`

    }

    //GETTING ALL DATA
    result := databases.DB.Raw("SELECT id_pajak, id_warga, tahun, total_tagihan_pajak FROM pajak WHERE id_warga = ? AND status_bayar = 'pending'", id).Scan(&tagihanPajak)

    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    if len(tagihanPajak) == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ada tagihan pajak tersedia"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"statusCode": http.StatusOK, "data": tagihanPajak})
}

func BayarPajakWarga(c *gin.Context) {
    //ID TAGIHAN
    id := c.Param("id")

    //VARIABLE FOR QUERY
    var resultStruct struct {
	IdWarga string `gorm:"column:id_warga"`
	TotalTagihanPajak string `gorm:"column:total_tagihan_pajak"`
    }

    //QUERY ID WARGA AND TOTAL_TAGIHAN_PAJAK
    result := databases.DB.Raw("SELECT id_warga, total_tagihan_pajak FROM pajak WHERE id_pajak = ?", id).Scan(&resultStruct)

    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    //CONVERT TO INTEGER
    totalTagihanPajakInt, err := strconv.Atoi(resultStruct.TotalTagihanPajak)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error converting totalTagihanPajak to int"})
        return
    }

    //CHECKING SALDO
    var saldo string
    resultSaldo := databases.DB.Raw("SELECT saldo FROM warga WHERE id_warga = ?", resultStruct.IdWarga).Row().Scan(&saldo)
    if resultSaldo != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": resultSaldo.Error()})
        return
    }

    //CONVERT SALDO TO INTEGER
    saldoInt, err := strconv.Atoi(saldo)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error converting saldo to int"})
        return
    }

    //CONDITIONS FOR PAYMENT
    if totalTagihanPajakInt > saldoInt {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Saldo tidak mencukupi",
                                            "statusCode": http.StatusBadRequest,
                                            "paymentStatus": "failed"})
        return
    } else {
        result := databases.DB.Exec("UPDATE pajak SET status_bayar = 'paid' WHERE id_pajak = ?", id)
        if result.Error != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
            return
        }

        result = databases.DB.Exec("UPDATE warga SET saldo = saldo - ? WHERE id_warga = ?", totalTagihanPajakInt, resultStruct.IdWarga)
        if result.Error != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"message": "Payment successfully processed", "statusCode": http.StatusOK, "paymentStatus": "success"})
            }

}
