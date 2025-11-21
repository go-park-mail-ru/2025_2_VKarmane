package image

type UploadImageResponse struct {
	ImageID string `json:"image_id"`
	URL     string `json:"url"`
}

type ImageURLResponse struct {
	URL string `json:"url"`
}
