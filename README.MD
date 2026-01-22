# Data Mahasiswa

Aplikasi web sederhana untuk mengelola data mahasiswa menggunakan Go dan Gin framework dengan database PostgreSQL.

## Fitur

- ✅ **Tambah Mahasiswa**: Menambah data mahasiswa baru (NIM, Nama, Nilai)
- ✅ **Validasi NIM**: NIM yang sudah terdaftar tidak bisa didaftarkan lagi
- ✅ **Lihat Data**: Menampilkan daftar semua mahasiswa dalam tabel
- ✅ **Edit Mahasiswa**: Mengedit data mahasiswa yang sudah ada
- ✅ **Hapus Mahasiswa**: Menghapus data mahasiswa dengan konfirmasi
- ✅ **Dark/Light Mode**: tampilan tema sesuai preferensi
- ✅ **Responsive Design**: Tampilan yang responsif di berbagai perangkat
- ✅ **Dua bahasa**: Tersedia bahasa Inggris dan bahasa Indonesia
- ✅ **Cari Mahasiswa**: Mencari data mahasiswa berdasarkan NIM atau Nama

## Prasyarat

Sebelum menjalankan aplikasi ini, pastikan Anda memiliki:

- [Go](https://golang.org/dl/) (versi 1.19 atau lebih baru)
- [PostgreSQL](https://www.postgresql.org/download/) (versi 12 atau lebih baru)
- Browser web modern

## Instalasi

1. **Clone atau download proyek ini**

2. **Install dependencies Go:**
   ```bash
   go mod tidy
   ```

3. **Setup Database PostgreSQL:**
   - Buat database baru bernama `kampus`
   - Buat tabel `mahasiswa` dengan struktur:
     ```sql
     CREATE TABLE mahasiswa (
         id SERIAL PRIMARY KEY,
         nim VARCHAR(20) NOT NULL,
         nama VARCHAR(100) NOT NULL,
         nilai DECIMAL(5,2) NOT NULL
     );
     ```

4. **Konfigurasi Database:**
   - Edit file `config.go` jika perlu mengubah konfigurasi database
   - Default: host=localhost, port=5432, user=postgres, password=postgres, dbname=kampus

## Menjalankan Aplikasi

1. **Jalankan server:**
   ```bash
   go run .
   ```

2. **Akses aplikasi:**
   - Buka browser dan kunjungi: `http://localhost:8080`

## Penggunaan

### Halaman Utama
- **Tambah Mahasiswa**: Isi form dan klik "Tambah Mahasiswa"
- **Lihat Data**: Data mahasiswa ditampilkan dalam tabel
- **Edit**: Klik "Edit" pada kolom Aksi untuk mengedit data
- **Hapus**: Klik "Hapus" pada kolom Aksi untuk menghapus data (dengan konfirmasi)
- **Toggle Tema**: Klik tombol 🌙 di kanan atas untuk mengganti dark/light mode
- **Bahasa**: Klik tombol ID/EN di kanan atas untuk mengganti bahasa
- **Pencarian**: Masukkan NIM atau Nama dan klik "Cari" untuk mencari data

### Halaman Edit
- Form edit akan terbuka dengan data mahasiswa sudah terisi
- Ubah data sesuai kebutuhan
- Klik "Update Mahasiswa" untuk menyimpan perubahan

## Teknologi yang Digunakan

- **Backend**: Go dengan Gin web framework
- **Database**: PostgreSQL dengan driver pgx
- **Frontend**: HTML, CSS, JavaScript
- **Template Engine**: Gin HTML templates

## Struktur Proyek

```
data-mahasiswa/
├── main.go           # Entry point aplikasi
├── models.go         # Definisi struct Mahasiswa
├── config.go         # Konfigurasi database
├── database.go       # Fungsi-fungsi database
├── handlers.go       # Handler HTTP untuk routes
├── templates/
│   ├── index.html    # Halaman utama
│   └── edit.html     # Halaman edit mahasiswa
├── go.mod            # Go modules
├── go.sum            # Go modules checksum
└── README.MD         # Dokumentasi ini
```

## API Endpoints

- `GET /` - Halaman utama (daftar mahasiswa)
- `POST /add` - Tambah mahasiswa baru
- `GET /edit/:id` - Halaman edit mahasiswa
- `POST /edit` - Update data mahasiswa
- `POST /delete` - Hapus mahasiswa

## Troubleshooting

### Error koneksi database
- Pastikan PostgreSQL sedang berjalan
- Periksa konfigurasi database di `config.go`
- Pastikan database `kampus` dan tabel `mahasiswa` sudah dibuat

### Port 8080 sudah digunakan
- Ganti port di `main.go`: `r.Run(":8081")`
- Atau hentikan proses yang menggunakan port 8080

### Template tidak dimuat
- Pastikan folder `templates/` ada dan file HTML di dalamnya valid
- Restart server setelah mengubah template

## Kontribusi

Untuk berkontribusi pada proyek ini:

1. Fork repository
2. Buat branch fitur baru (`git checkout -b feature/AmazingFeature`)
3. Commit perubahan (`git commit -m 'Add some AmazingFeature'`)
4. Push ke branch (`git push origin feature/AmazingFeature`)
5. Buat Pull Request

## Lisensi

Proyek ini menggunakan lisensi MIT. Lihat file `LICENSE` untuk detail lebih lanjut.

---

**Catatan**: Ini adalah proyek pembelajaran untuk demonstrasi CRUD operations dengan Go dan PostgreSQL.