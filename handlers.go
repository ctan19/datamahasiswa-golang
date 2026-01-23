package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func indexHandler(c *gin.Context, db *sql.DB) {
	// Ambil parameter sorting dari query string
	sortBy := c.DefaultQuery("sort_by", "id") // Default: sort by ID
	order := c.DefaultQuery("order", "asc")   // Default: ascending

	// Ambil parameter pencarian
	searchQuery := c.Query("search")
	minNilaiStr := c.Query("min_nilai")
	var minNilai *float64
	if minNilaiStr != "" {
		parsedNilai, err := strconv.ParseFloat(minNilaiStr, 64)
		if err != nil {
			mahasiswa, _ := getMahasiswa(db)
			c.HTML(http.StatusBadRequest, "index.html", gin.H{"error": "Nilai harus angka", "mahasiswa": mahasiswa})
			return
		}
		minNilai = &parsedNilai
	}

	var mahasiswa []Mahasiswa
	var err error

	if searchQuery != "" || minNilai != nil {
		// Jika ada query pencarian, gunakan fungsi pencarian
		mahasiswa, err = searchMahasiswaAdvanced(db, searchQuery, minNilai)
	} else {
		// Jika tidak ada pencarian, gunakan sorting biasa
		mahasiswa, err = getMahasiswaSorted(db, sortBy, order)
	}

	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{"error": err.Error(), "mahasiswa": []Mahasiswa{}, "total": 0})
		return
	}

	// Cek query parameter untuk pesan sukses
	var deletedNIM string
	if deleted := c.Query("deleted"); deleted == "true" {
		if nim := c.Query("nim"); nim != "" {
			deletedNIM = nim
		}
	}

	total := len(mahasiswa)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"mahasiswa":   mahasiswa,
		"deleted_nim": deletedNIM,
		"total":       total,
		"sort_by":     sortBy, // Kirim info sorting ke template
		"order":       order,
		"search_text": searchQuery,
		"min_nilai":   minNilaiStr,
	})
}

func addHandler(c *gin.Context, db *sql.DB) {
	nama := c.PostForm("nama")
	nim := c.PostForm("nim")
	nilaiStr := c.PostForm("nilai")
	nilai, err := strconv.ParseFloat(nilaiStr, 64)
	if err != nil {
		// Jika error parsing nilai, tetap tampilkan data mahasiswa yang ada
		mahasiswa, _ := getMahasiswa(db)
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"error": "Nilai harus angka", "mahasiswa": mahasiswa})
		return
	}
	if nilai > 4 {
		mahasiswa, _ := getMahasiswa(db)
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"error": "Nilai IPK maksimal 4", "mahasiswa": mahasiswa})
		return
	}

	err = tambahMahasiswa(db, nama, nim, nilai)
	if err != nil {
		// Jika error (seperti NIM sudah ada), tetap tampilkan data mahasiswa yang ada
		mahasiswa, _ := getMahasiswa(db)
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"error": err.Error(), "mahasiswa": mahasiswa})
		return
	}

	c.Redirect(http.StatusFound, "/")
}

func deleteHandler(c *gin.Context, db *sql.DB) {
	idStr := c.PostForm("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Jika ID tidak valid, tetap tampilkan data mahasiswa yang ada
		mahasiswa, _ := getMahasiswa(db)
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"error": "ID tidak valid", "mahasiswa": mahasiswa})
		return
	}

	// Ambil data mahasiswa sebelum dihapus untuk mendapatkan NIM
	mahasiswa, err := getMahasiswaByID(db, id)
	if err != nil {
		// Jika mahasiswa tidak ditemukan, tetap tampilkan data mahasiswa yang ada
		allMahasiswa, _ := getMahasiswa(db)
		c.HTML(http.StatusNotFound, "index.html", gin.H{"error": "Mahasiswa tidak ditemukan", "mahasiswa": allMahasiswa})
		return
	}

	err = deleteMahasiswa(db, id)
	if err != nil {
		// Jika error delete, tetap tampilkan data mahasiswa yang ada
		allMahasiswa, _ := getMahasiswa(db)
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{"error": err.Error(), "mahasiswa": allMahasiswa})
		return
	}

	// Redirect dengan pesan sukses
	c.Redirect(http.StatusFound, "/?deleted=true&nim="+mahasiswa.NIM)
}

func editPageHandler(c *gin.Context, db *sql.DB) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"error": "ID tidak valid", "mahasiswa": []Mahasiswa{}})
		return
	}

	mahasiswa, err := getMahasiswaByID(db, id)
	if err != nil {
		// Jika mahasiswa tidak ditemukan, tetap tampilkan data mahasiswa yang ada
		allMahasiswa, _ := getMahasiswa(db)
		c.HTML(http.StatusNotFound, "index.html", gin.H{"error": "Mahasiswa tidak ditemukan", "mahasiswa": allMahasiswa})
		return
	}

	// Cek query parameter untuk pesan sukses
	success := c.Query("success") == "true"

	c.HTML(http.StatusOK, "edit.html", gin.H{"mahasiswa": mahasiswa, "success": success})
}

func updateHandler(c *gin.Context, db *sql.DB) {
	idStr := c.PostForm("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"error": "ID tidak valid", "mahasiswa": []Mahasiswa{}})
		return
	}

	nama := c.PostForm("nama")
	nim := c.PostForm("nim")
	nilaiStr := c.PostForm("nilai")
	nilai, err := strconv.ParseFloat(nilaiStr, 64)
	if err != nil {
		// Jika error parsing nilai, tetap tampilkan data mahasiswa yang ada
		mahasiswa, _ := getMahasiswa(db)
		c.HTML(http.StatusBadRequest, "index.html", gin.H{"error": "Nilai harus angka", "mahasiswa": mahasiswa})
		return
	}
	if nilai > 4 {
		c.HTML(http.StatusBadRequest, "edit.html", gin.H{"error": "Nilai IPK maksimal 4", "mahasiswa": &Mahasiswa{ID: id, Nama: nama, NIM: nim, Nilai: nilai}})
		return
	}

	err = updateMahasiswa(db, id, nama, nim, nilai)
	if err != nil {
		// Jika error, redirect ke halaman edit dengan pesan error
		c.HTML(http.StatusBadRequest, "edit.html", gin.H{"error": err.Error(), "mahasiswa": &Mahasiswa{ID: id, Nama: nama, NIM: nim, Nilai: nilai}})
		return
	}

	// Ambil data mahasiswa yang telah diupdate
	mahasiswa, err := getMahasiswaByID(db, id)
	if err != nil {
		// Jika gagal ambil data, redirect ke halaman utama
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.HTML(http.StatusOK, "edit.html", gin.H{"mahasiswa": mahasiswa, "success": true})
}
