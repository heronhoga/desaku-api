package warga

import (
	"github.com/gin-gonic/gin"
	"desaku-api/databases"
	"desaku-api/models"
	"net/http"
	"strconv"
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

	//REQUEST BODY (id_warga)
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

	//CHECK IF THERE'S AN EXISTING USER
	var existingWarga string
	querySearchWarga := databases.DB.Raw("SELECT id_warga FROM warga WHERE id_warga = ?", wifiData.IdWarga).Scan(&existingWarga)

	if querySearchWarga.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Warga tidak valid"})
		return
	} else if querySearchWarga.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": querySearchWarga.Error.Error()})
		return
	}

	//CHECK IF USER IS ALREADY SUBSCRIBED
	var langgananWifi models.LanggananWifi
	result := databases.DB.Raw("SELECT * FROM daftar_pelanggan_wifi WHERE id_warga = ?", wifiData.IdWarga).Scan(&langgananWifi)

	if result.RowsAffected > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "User already subscribed"})
		return
	}

	//INSERT DATA / NEW SUBSCRIPTION
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

	//VARIABLE FOR REQUEST BODY (id_warga)
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

	//CHECK IF THERE IS A SUBSCRIPTION
	var existingWarga string
	querySearchWarga := databases.DB.Raw("SELECT id_warga FROM daftar_pelanggan_wifi WHERE id_warga = ?", wifiData.IdWarga).Scan(&existingWarga)

	if querySearchWarga.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Warga tidak valid"})
		return
	} else if querySearchWarga.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": querySearchWarga.Error.Error()})
		return
	}

	//UNSUBSCRIBE
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

func BayarWifiWarga(c *gin.Context) {
    id := c.Param("id")
    
	//VARIABLES FOR TAGIHAN
    var tagihan struct {
		//ALIAS ID WARGA
		IdPelanggan string `gorm:"column:id_pelanggan"`
		//ID TAGIHAN
        TotalTagihanWifi string `gorm:"column:total_tagihan_wifi"`
    }

	//CARI TAGIHAN
    result := databases.DB.Raw("SELECT id_pelanggan, total_tagihan_wifi FROM tagihan_wifi WHERE id_tagihan = ? AND status = 'pending'", id).Scan(&tagihan)

    if result.RowsAffected == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Tidak ada tagihan tersedia"})
        return
    } else if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

	//CONVERT TAGIHAN TO INT
	totalTagihanInt, err := strconv.Atoi(tagihan.TotalTagihanWifi)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error converting total_tagihan to int"})
        return
    }

	//VARIABLE FOR SALDO WARGA
    var warga struct {
        Saldo string `gorm:"column:saldo"`
    }

	//CARI SALDO
    result = databases.DB.Raw("SELECT saldo FROM warga WHERE id_warga = ?", tagihan.IdPelanggan).Scan(&warga)

    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }
	
	//CONVERT SALDO TO INT
	totalSaldoInt, err := strconv.Atoi(warga.Saldo)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error converting saldo to int"})
        return
    }

	//CONDITIONS
    if totalSaldoInt >= totalTagihanInt {
		//UPDATE PAYMENT
        sql := "UPDATE tagihan_wifi SET status = ? WHERE id_tagihan = ?"
        result = databases.DB.Exec(sql, "paid", id)
        if result.Error != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
            return
        }

		//KURANGI SALDO
		sql = "UPDATE warga SET saldo = ? WHERE id_warga = ?"
		result = databases.DB.Exec(sql, strconv.Itoa(totalSaldoInt - totalTagihanInt), tagihan.IdPelanggan)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

        c.JSON(http.StatusOK, gin.H{
            "statusCode": http.StatusOK,
            "message": "Payment successful",
        })
    } else {
        c.JSON(http.StatusForbidden, gin.H{
            "statusCode": http.StatusForbidden,
            "message": "Saldo tidak cukup",
        })
    }
}
