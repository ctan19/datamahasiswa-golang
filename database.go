package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func OpenDB(cfg DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxIdleTime(30 * time.Minute)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}

func getMahasiswa(db *sql.DB) ([]Mahasiswa, error) {
	rows, err := db.Query(`SELECT id, nim, nama, nilai FROM mahasiswa ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mahasiswa []Mahasiswa
	for rows.Next() {
		var m Mahasiswa
		if err := rows.Scan(&m.ID, &m.NIM, &m.Nama, &m.Nilai); err != nil {
			return nil, err
		}
		mahasiswa = append(mahasiswa, m)
	}
	return mahasiswa, nil
}

func getMahasiswaSorted(db *sql.DB, sortBy, order string) ([]Mahasiswa, error) {
	// Validasi parameter untuk mencegah SQL injection
	validSortBy := map[string]bool{"id": true, "nim": true, "nama": true, "nilai": true}
	validOrder := map[string]bool{"asc": true, "desc": true}

	if !validSortBy[sortBy] {
		sortBy = "id"
	}
	if !validOrder[order] {
		order = "asc"
	}

	query := fmt.Sprintf(`SELECT id, nim, nama, nilai FROM mahasiswa ORDER BY %s %s`, sortBy, order)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mahasiswa []Mahasiswa
	for rows.Next() {
		var m Mahasiswa
		if err := rows.Scan(&m.ID, &m.NIM, &m.Nama, &m.Nilai); err != nil {
			return nil, err
		}
		mahasiswa = append(mahasiswa, m)
	}
	return mahasiswa, nil
}

func checkNIMExists(db *sql.DB, nim string) (bool, error) {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM mahasiswa WHERE nim = $1`, nim).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func checkNIMExistsExceptID(db *sql.DB, nim string, excludeID int) (bool, error) {
	var count int
	err := db.QueryRow(`SELECT COUNT(*) FROM mahasiswa WHERE nim = $1 AND id != $2`, nim, excludeID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func tambahMahasiswa(db *sql.DB, nama, nim string, nilai float64) error {
	exists, err := checkNIMExists(db, nim)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("NIM sudah terdaftar")
	}

	_, err = db.Exec(
		`INSERT INTO mahasiswa (nim, nama, nilai) VALUES ($1,$2,$3)`,
		nim, nama, nilai,
	)
	return err
}

func deleteMahasiswa(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM mahasiswa WHERE id = $1`, id)
	return err
}

func updateMahasiswa(db *sql.DB, id int, nama, nim string, nilai float64) error {
	exists, err := checkNIMExistsExceptID(db, nim, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("NIM sudah terdaftar oleh mahasiswa lain")
	}

	_, err = db.Exec(`UPDATE mahasiswa SET nim = $1, nama = $2, nilai = $3 WHERE id = $4`, nim, nama, nilai, id)
	return err
}

func getMahasiswaByID(db *sql.DB, id int) (*Mahasiswa, error) {
	var m Mahasiswa
	err := db.QueryRow(`SELECT id, nim, nama, nilai FROM mahasiswa WHERE id = $1`, id).Scan(&m.ID, &m.NIM, &m.Nama, &m.Nilai)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func searchMahasiswa(db *sql.DB, query string) ([]Mahasiswa, error) {
	rows, err := db.Query(`SELECT id, nim, nama, nilai FROM mahasiswa WHERE nim = $1 OR nama ILIKE $2 ORDER BY id`, query, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var mahasiswa []Mahasiswa
	for rows.Next() {
		var m Mahasiswa
		if err := rows.Scan(&m.ID, &m.NIM, &m.Nama, &m.Nilai); err != nil {
			return nil, err
		}
		mahasiswa = append(mahasiswa, m)
	}
	return mahasiswa, nil
}
