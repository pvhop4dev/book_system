package dto

type FileResponse struct {
	FileName string `json:"file_name"`
	URL      string `json:"url,omitempty"`
	Size     int64  `json:"size,omitempty"`
}
