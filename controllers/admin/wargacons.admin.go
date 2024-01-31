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

func GetSpecificWarga(c *gin.Context) {
	id := c.Param("id")

	var warga struct {
		IdWarga string `json:"id_warga"`
		Nama string `json:"nama"`
		TanggalLahir string `json:"tanggal_lahir"`
		JenisKelamin string `json:"jenis_kelamin"`
		Nik string `json:"nik"`
		Alamat string `json:"alamat"`
		Saldo string `json:"saldo"`
	}

	wargaQuery := databases.DB.Raw(`SELECT id_warga, nama, tanggal_lahir,jenis_kelamin, nik, alamat, saldo FROM warga WHERE id_warga = ?`, id).Scan(&warga)
	if wargaQuery.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": wargaQuery.Error.Error()})
		return
	}

	if wargaQuery.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No record found"})
	}

	c.JSON(http.StatusOK, gin.H{"data": warga})
}

func EditWargaData(c *gin.Context) {
    id := c.Param("id")

    var requestData struct {
        Nama         *string `json:"nama" binding:"-"`
        TanggalLahir *string `json:"tanggal_lahir" binding:"-"`
        JenisKelamin *string `json:"jenis_kelamin" binding:"-"`
        Nik          *string `json:"nik" binding:"-"`
        Alamat       *string `json:"alamat" binding:"-"`
        Saldo        *string `json:"saldo" binding:"-"`
    }
	
	var previousData struct {
        Nama         *string `json:"nama"`
        TanggalLahir *string `json:"tanggal_lahir"`
        JenisKelamin *string `json:"jenis_kelamin"`
        Nik          *string `json:"nik"`
        Alamat       *string `json:"alamat"`
        Saldo        *string `json:"saldo"`
    }

    if err := c.ShouldBindJSON(&requestData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

	//GET PREVIOUS DATA
	previousQuery := databases.DB.Raw(`SELECT nama, tanggal_lahir, jenis_kelamin, nik, alamat, saldo FROM warga WHERE id_warga = ?`, id).Scan(&previousData)
	if previousQuery.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": previousQuery.Error.Error()})
		return
	}

	if requestData.Nama == nil {
	    requestData.Nama = previousData.Nama
	}
	if requestData.TanggalLahir == nil {
	    requestData.TanggalLahir = previousData.TanggalLahir
	}
	if requestData.JenisKelamin == nil {
	    requestData.JenisKelamin = previousData.JenisKelamin
	}
	if requestData.Nik == nil {
	    requestData.Nik = previousData.Nik
	}
	if requestData.Alamat == nil {
	    requestData.Alamat = previousData.Alamat
	}
	if requestData.Saldo == nil {
	    requestData.Saldo = previousData.Saldo
	}

	updateQuery := databases.DB.Exec("UPDATE warga SET nama = ?, tanggal_lahir = ?, jenis_kelamin = ?, nik = ?, alamat = ?, saldo = ? WHERE id_warga = ?", requestData.Nama, requestData.TanggalLahir, requestData.JenisKelamin, requestData.Nik, requestData.Alamat, requestData.Saldo, id)

	if updateQuery.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": updateQuery.Error.Error()})
		return
	}

	if updateQuery.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No record found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data updated successfully"})
}
