package dto

type FileDTO struct {
	FileName   string `json:"file_name"`
	FileType   string `json:"file_type"`
	FileSize   int64  `json:"file_size"`
	UploadTime string `json:"upload_time"`
}
