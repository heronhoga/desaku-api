package warga

import (
	"desaku-api/databases"
	"desaku-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//REGISTER --WARGA--
func RegisterWarga(c *gin.Context) {
    var warga models.Warga

    if err := c.ShouldBindJSON(&warga); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if warga.Nama == "" || warga.TanggalLahir == "" || warga.JenisKelamin == "" || warga.Nik == "" || warga.Alamat == "" || warga.Password == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
        return
    }

    var existingWarga models.Warga
	result := databases.DB.Raw("SELECT nama FROM warga WHERE nama = ?", warga.Nama).Scan(&existingWarga)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
    c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
    return
	}
	
    if existingWarga.Nama != "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Name already exists"})
        return
    }

    sql := "INSERT INTO warga(nama, tanggal_lahir, jenis_kelamin, nik, alamat, password) VALUES(?, ?, ?, ?, ?, ?)"
    result = databases.DB.Exec(sql, warga.Nama, warga.TanggalLahir, warga.JenisKelamin, warga.Nik, warga.Alamat, warga.Password)

    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "statusCode": http.StatusOK,
        "data": warga})
}

	
	