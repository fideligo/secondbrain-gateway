package model

type DocumentRequestJSON struct {
	FileName string `json:"file_name"`
	Author   string `json:"author"`
}