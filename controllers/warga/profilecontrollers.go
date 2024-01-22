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
    //REQUEST BODY
    var warga models.Warga

    if err := c.ShouldBindJSON(&warga); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    //CHECK IF ALL FIELDS ARE FILLED
    if warga.Nama == "" || warga.TanggalLahir == "" || warga.JenisKelamin == "" || warga.Nik == "" || warga.Alamat == "" || warga.Password == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
        return
    }

    //CHECK IF DATA ALREADY EXISTS
    var existingWarga models.Warga
    result := databases.DB.Raw("SELECT nama FROM warga WHERE nama = ? OR nik = ?", warga.Nama, warga.Nik).Scan(&existingWarga)
    if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    if existingWarga.Nama != "" || existingWarga.Nik != "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Data already exists"})
        return
    }

    //INSERT DATA
    query := "INSERT INTO warga (nama, tanggal_lahir, jenis_kelamin, nik, alamat, password) VALUES (?,?,?,?,?,?)"
    resultInsert := databases.DB.Exec(query, warga.Nama, warga.TanggalLahir, warga.JenisKelamin, warga.Nik, warga.Alamat, warga.Password)

    if resultInsert.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    //GET ID
    queryId := "SELECT id_warga FROM warga WHERE nama = ?"
    var idWarga int
    row := databases.DB.Raw(queryId, warga.Nama).Row()
    err := row.Scan(&idWarga)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "statusCode": http.StatusOK,
        "data": warga,
        "id_warga": idWarga,})
}

//LOGIN WARGA
func LoginWarga(c *gin.Context) {
    //REQUEST BODY
    var loginData struct {
    Nama     string `json:"nama"`
    Password string `json:"password"`
}

    if err := c.ShouldBindJSON(&loginData); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    //CHECK IF ALL FIELDS ARE FILLED
    if loginData.Nama == "" || loginData.Password == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Both 'nama' and 'password' fields must be filled"})
        return
    }

    //MATCHING DATA
    var warga models.Warga
    result := databases.DB.Raw("SELECT * FROM warga WHERE nama = ? AND password = ?", loginData.Nama, loginData.Password).Scan(&warga)

    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    if warga.Nama == "" || warga.Password == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid username or password",
            "loginStatus": false,})
        return
    }

    //GET ID
    var idWarga int64
    row := databases.DB.Raw("SELECT id_warga FROM warga WHERE nama = ? AND password = ?", loginData.Nama, loginData.Password).Row()
    err := row.Scan(&idWarga)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "statusCode": http.StatusOK,
        "loginStatus": true,
        "id_warga": idWarga,
        "data": warga})
}

//GET PROFILE DATA WARGA
func ProfileWargaOne(c *gin.Context) {
    //ID PARAMETER
    id := c.Param("id")
    var warga models.Warga
    result := databases.DB.Raw("SELECT * FROM warga WHERE id_warga = ?", id).Scan(&warga)

    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
        return
    }

    //CHECK IF USER IS NOT EXISTS
    if warga.Nama == "" || warga.Nik == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "No profile found with given id",
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "statusCode": http.StatusOK,
        "id_warga": id,
        "data": warga})
}
