package database

import (
	"database/sql"
	"log"

	_"github.com/lib/pq"
)

func InitDB(connectionString string) (*sql.DB, error) {
	// Buka database
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	// Tes koneksi
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Set pengaturan koneksi (opsional namun direkomendasikan)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Database sukses terhubung")
	return db, nil
}