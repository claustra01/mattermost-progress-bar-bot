package bot

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/claustra01/mattermost-progress-bar-bot/pkg/date"
)

func PostMessage(baseUrl string, channelID string, token string, fileKey string) {
	// create message
	now := time.Now()
	progress := date.GetProgress(now)
	remainingDays := date.GetRemainingDays(now)

	var message string
	if remainingDays > 0 {
		message = fmt.Sprintf("プログラム全体の%.f%%が経過しました。成果発表会まであと%d日です。", progress*100, remainingDays)
	} else if remainingDays == 0 {
		message = "プログラム全体の100%が経過しました。本日は成果発表会です。"
	}

	// create request body
	body := CreatePostRequestBody{
		ChannelID: channelID,
		Message:   message,
		RootID:    nil,
		FileIDs:   []string{fileKey},
	}
	data, err := MarshalCreatePostReqBody(body)
	if err != nil {
		slog.Error("Error marshalling request body:", err)
		return
	}

	// create request
	url := fmt.Sprintf("%s/api/v4/posts", baseUrl)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		slog.Error("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// send request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading response:", err)
		return
	}
	slog.Info("Reponse:", "URL", url, "Body", string(respBody))

	// shutdown bot
	if remainingDays <= 0 {
		os.Exit(0)
	}
}

func UploadImage(baseUrl string, channelID string, token string, filename string) string {
	// open image file
	file, err := os.Open(filename)
	if err != nil {
		slog.Error("Error opening file:", err)
		return ""
	}
	defer file.Close()

	// create form
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	w.WriteField("channel_id", channelID)
	part, err := w.CreateFormFile("files", filepath.Base(filename))
	if err != nil {
		slog.Error("Error creating form file:", err)
		return ""
	}
	_, err = io.Copy(part, file)
	if err != nil {
		slog.Error("Error copying file to form file:", err)
		return ""
	}
	if err := w.Close(); err != nil {
		slog.Error("Error closing writer:", err)
		return ""
	}

	// create request
	url := fmt.Sprintf("%s/api/v4/files", baseUrl)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		slog.Error("Error creating request:", err)
		return ""
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", w.FormDataContentType())

	// send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Error sending request:", err)
		return ""
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading response:", err)
		return ""
	}

	data, err := UnmarshalUploadFileResponseBody(respBody)
	if err != nil {
		slog.Error("Error unmarshalling response body:", err)
		return ""
	}

	return data.FileInfos[0].ID
}
