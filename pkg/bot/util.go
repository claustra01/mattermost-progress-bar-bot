package bot

import (
	"encoding/json"
)

type CreatePostRequestBody struct {
	ChannelID string   `json:"channel_id"`
	Message   string   `json:"message"`
	RootID    *string  `json:"root_id"`
	FileIDs   []string `json:"file_ids"`
}

func MarshalCreatePostReqBody(body CreatePostRequestBody) ([]byte, error) {
	raw, err := json.Marshal(body)
	if err != nil {
		return raw, err
	}
	return raw, nil
}

type UploadedFileInfo struct {
	ID              string `json:"id"`
	UserID          string `json:"user_id"`
	PostID          string `json:"post_id"`
	CreatedAt       int    `json:"created_at"`
	UpdatedAt       int    `json:"updated_at"`
	DeleteAt        int    `json:"delete_at"`
	Name            string `json:"name"`
	Extension       string `json:"extension"`
	Size            int    `json:"size"`
	MimeType        string `json:"mime_type"`
	Width           int    `json:"width"`
	Height          int    `json:"height"`
	HasPreviewImage bool   `json:"has_preview_image"`
}

type UploadFileResponseBody struct {
	FileInfos []UploadedFileInfo `json:"file_infos"`
	ClientIDs []string           `json:"client_ids"`
}

func UnmarshalUploadFileResponseBody(raw []byte) (UploadFileResponseBody, error) {
	var res UploadFileResponseBody
	err := json.Unmarshal(raw, &res)
	if err != nil {
		return res, err
	}
	return res, nil
}
