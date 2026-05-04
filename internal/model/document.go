package model

import (
	"time"
	"gorm.io/gorm"
)

// DocumentRequestJSON untuk menangkap request dari API
type DocumentRequestJSON struct {
	FileName string `json:"file_name"`
	Author   string `json:"author"`
}

// Document adalah struktur untuk tabel di PostgreSQL
type Document struct {
	gorm.Model // Ini otomatis menambahkan ID, CreatedAt, UpdatedAt, DeletedAt
	FileName   string    `json:"file_name"`
	Author     string    `json:"author"`
	FilePath   string    `json:"file_path"`
	Summary    string    `json:"summary"`
	UploadedAt time.Time `json:"uploaded_at"`
}