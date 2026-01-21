package main

type Mahasiswa struct {
	ID    int     `json:"id"`
	Nama  string  `json:"nama"`
	NIM   string  `json:"nim"`
	Nilai float64 `json:"nilai"`
}

type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}
