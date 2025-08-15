package bot

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/claustra01/typetalk-progress-bar-bot/pkg/date"
)

func UploadImage(topicId string, token string, filename string) string {
	url := fmt.Sprintf("https://typetalk.com/api/v1/topics/%s/attachments?typetalkToken=%s", topicId, token)

	file, err := os.Open(filename)
	if err != nil {
		slog.Error("Error opening file:", err)
		return ""
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		slog.Error("Error creating form file:", err)
		return ""
	}

	_, err = io.Copy(part, file)
	if err != nil {
		slog.Error("Error copying file to form file:", err)
		return ""
	}

	err = writer.Close()
	if err != nil {
		slog.Error("Error closing writer:", err)
		return ""
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		slog.Error("Error creating request:", err)
		return ""
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Error sending request:", err)
		return ""
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Error reading response:", err)
		return ""
	}

	data, err := UnmarshalRespBody(responseBody)
	if err != nil {
		slog.Error("Error unmarshalling response body:", err)
		return ""
	}

	return data.FileKey
}

func PostMessage(topicId string, token string, fileKey string) {
	url := fmt.Sprintf("https://typetalk.com/api/v1/topics/%s?typetalkToken=%s", topicId, token)

	now := time.Now()
	progress := date.GetProgress(now)
	remainingDays := date.GetRemainingDays(now)

	var message string
	if remainingDays > 0 {
		message = fmt.Sprintf("プログラム全体の%.f%%が経過しました。成果発表会まであと%d日です。", progress*100, remainingDays)
	} else if remainingDays == 0 {
		message = "プログラム全体の100%が経過しました。本日は成果発表会です。"
	} else {
		os.Exit(1)
	}

	body := PostMessageRequest{
		Message:  message,
		FileKeys: []string{fileKey},
	}
	data, err := MarshalReqBody(body)
	if err != nil {
		slog.Error("Error marshalling request body:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		slog.Error("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
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

	slog.Info("Request:", "URL", url, "Body", string(respBody))
}
