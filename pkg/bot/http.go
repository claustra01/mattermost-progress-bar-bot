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

func PostMessage(baseUrl string, channelID string, token string, fileKey string) error {
	// create message
	now := time.Now()
	progress := date.GetProgress(now)
	remainingDays := date.GetRemainingDays(now)

	if remainingDays < 0 {
		return fmt.Errorf("invalid remaining days: %d", remainingDays)
	}

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
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	// create request
	url := fmt.Sprintf("%s/api/v4/posts", baseUrl)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// send request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("error response from server: StatusCode %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	slog.Info("Response:", "URL", url, "Body", string(respBody))
	return nil
}

func UploadImage(baseUrl string, channelID string, token string, filename string) (string, error) {
	// open image file
	file, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	// create form
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	w.WriteField("channel_id", channelID)
	part, err := w.CreateFormFile("files", filepath.Base(filename))
	if err != nil {
		return "", fmt.Errorf("failed to create form file for %s: %w", filename, err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("failed to copy file %s to form file: %w", filename, err)
	}
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	// create request
	url := fmt.Sprintf("%s/api/v4/files", baseUrl)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", w.FormDataContentType())

	// send request
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("error response from server: StatusCode %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %w", err)
	}

	data, err := UnmarshalUploadFileResponseBody(respBody)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling response body: %w", err)
	}

	return data.FileInfos[0].ID, nil
}
