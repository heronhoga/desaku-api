package models

type Warga struct {
    Nama          string `json:"nama"`
    TanggalLahir  string `json:"tanggal_lahir"`
    JenisKelamin  string `json:"jenis_kelamin"`
    Nik           string `json:"nik"`
    Alamat        string `json:"alamat"`
    Password      string `json:"password"`
}