package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// DB adalah koneksi global ke database
// bisa dipanggil dari file manapun
var DB *sql.DB

func Connect() {
	// Ambil konfigurasi dari file .env
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Format string koneksi PostgreSQL
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	// Buka koneksi
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Gagal konek ke database:", err)
	}

	// Cek apakah koneksi beneran berhasil
	if err = DB.Ping(); err != nil {
		log.Fatal("Database tidak bisa di-ping:", err)
	}

	log.Println("Berhasil konek ke database!")
}

func Migrate() {
	query := `
	CREATE TABLE IF NOT EXISTS labs (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		location VARCHAR(100),
		capacity INT DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS devices (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		ip_address VARCHAR(20) NOT NULL UNIQUE,
		lab_id INT REFERENCES labs(id),
		status VARCHAR(10) DEFAULT 'unknown',
		last_seen TIMESTAMP
	);`

	_, err := DB.Exec(query)
	if err != nil {
		log.Fatal("Gagal migrate tabel:", err)
	}

	log.Println("table siap")
}