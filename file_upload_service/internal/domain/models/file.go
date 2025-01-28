package models

type FileModel struct {
	ID         int64  `db:"id"`
	FileName   string `db:"file_name"`
	FileType   string `db:"file_type"`
	FileSize   int64  `db:"file_size"`
	UploadTime string `db:"upload_time"`
	PathToFile string `db:"path_to_file"`
}
