package bot

import (
	"encoding/json"
)

type UploadImageResponse struct {
	ContentType string `json:"contentType"`
	FileKey     string `json:"fileKey"`
	FileName    string `json:"fileName"`
	FileSize    int    `json:"fileSize"`
}

func UnmarshalRespBody(raw []byte) (UploadImageResponse, error) {
	var res UploadImageResponse
	err := json.Unmarshal(raw, &res)
	if err != nil {
		return res, err
	}
	return res, nil
}

type PostMessageRequest struct {
	Message  string   `json:"message"`
	FileKeys []string `json:"fileKeys"`
}

func MarshalReqBody(body PostMessageRequest) ([]byte, error) {
	raw, err := json.Marshal(body)
	if err != nil {
		return raw, err
	}
	return raw, nil
}
