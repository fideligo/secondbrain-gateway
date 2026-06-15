package model

import (
	"time"
	"gorm.io/gorm"
)

type DocumentRequestJSON struct {
	FileName string `json:"file_name"`
	Author   string `json:"author"`
}

type Document struct {
	gorm.Model
	FileName   string    `json:"file_name"`
	Author     string    `json:"author"`
	FilePath   string    `json:"file_path"`
	Summary    string    `json:"summary"`
	UploadedAt time.Time `json:"uploaded_at"`
}