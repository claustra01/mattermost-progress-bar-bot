package bot

import (
	"fmt"
	"image"
	"image/color"
	"log/slog"
	"time"

	"github.com/claustra01/typetalk-progress-bar-bot/pkg/date"
	"gocv.io/x/gocv"
)

func GenerateImage() string {
	// 画像のサイズを設定
	width := 1000
	height := 200

	// 画像を作成
	img := gocv.NewMatWithSize(height, width, gocv.MatTypeCV8UC3)
	defer img.Close()

	// 背景を塗りつぶし
	bgColor := color.RGBA{255, 255, 255, 255}
	gocv.Rectangle(&img, image.Rect(0, 0, width, height), bgColor, -1)

	// プログレスバーを描画
	now := time.Now()
	progress := date.GetProgress(now)
	progressBarColor := color.RGBA{255, 0, 0, 255}
	progressWidth := int(float64(width-40) * progress)
	gocv.Rectangle(&img, image.Rect(20, 60, 20+progressWidth, height-20), progressBarColor, -1)

	// プログレスバーの枠を描画
	borderColor := color.RGBA{32, 32, 32, 255}
	gocv.Rectangle(&img, image.Rect(20, 60, width-20, height-20), borderColor, 10)

	// タイトルを描画
	title := "SecHack365 Progress Bar"
	gocv.PutText(&img, title, image.Point{width/2 - 220, 40}, gocv.FontHersheyComplex, 1.0, borderColor, 2)

	// 画像を保存
	filename := fmt.Sprintf("progress_bar_%d_%02d_%02d.png", now.Year(), now.Month(), now.Day())
	gocv.IMWrite(filename, img)
	slog.Info("Generated image:", filename)

	return filename
}
