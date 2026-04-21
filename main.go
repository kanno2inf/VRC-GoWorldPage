package main

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const VRCHATWorldURLFormat = "https://vrchat.com/home/world/%s/info"

// 画像ファイルからXMPメタデータ文字列を抽出
func loadXMP(imagePath string) (string, error) {
	data, err := os.ReadFile(imagePath)
	if err != nil {
		return "", err
	}

	startTag := []byte("<x:xmpmeta")
	endTag := []byte("</x:xmpmeta>")

	start := bytes.Index(data, startTag)
	if start == -1 {
		return "", nil
	}

	end := bytes.Index(data[start:], endTag)
	if end == -1 {
		return "", nil
	}
	end += start + len(endTag)

	return string(data[start:end]), nil
}

// 画像ファイルのXMPメタデータからWorldIDを抽出
func getWorldIDFromImage(imagePath string) (string, error) {
	xmpString, err := loadXMP(imagePath)
	if err != nil {
		return "", err
	}
	if xmpString == "" {
		return "", nil
	}

	decoder := xml.NewDecoder(bytes.NewReader([]byte(xmpString)))

	for {
		tok, err := decoder.Token()
		if err != nil {
			break
		}

		switch se := tok.(type) {
		case xml.StartElement:
			if strings.HasSuffix(se.Name.Local, "WorldID") {
				var worldID string
				if err := decoder.DecodeElement(&worldID, &se); err != nil {
					return "", err
				}
				return worldID, nil
			}
		}
	}

	return "", nil
}

// デフォルトブラウザを開く
func openWorldPage(worldID string) error {
	worldURL := fmt.Sprintf(VRCHATWorldURLFormat, worldID)

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", worldURL)
	case "darwin":
		cmd = exec.Command("open", worldURL)
	default:
		cmd = exec.Command("xdg-open", worldURL)
	}

	return cmd.Start()
}

// キー入力されるまで待機
func waitForEnter() {
	fmt.Println()
	fmt.Println("Enterキーを押すと終了します...")
	reader := bufio.NewReader(os.Stdin)
	_, _ = reader.ReadString('\n')
}

// pngファイルかどうかチェック
func isPNGFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".png"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("画像ファイルを指定してください.")
		waitForEnter()
		os.Exit(1)
	}

	paths := os.Args[1:]
	uniqueWorldIDs := make(map[string]struct{})
	hadError := false
	openedCount := 0

	for _, path := range paths {
		if !isPNGFile(path) {
			fmt.Fprintf(os.Stderr, "無効なファイル形式: %s\n", filepath.Base(path))
			hadError = true
			continue
		}

		worldID, err := getWorldIDFromImage(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "エラー: %s: %v\n", path, err)
			hadError = true
			continue
		}
		if worldID == "" {
			continue
		}
		uniqueWorldIDs[worldID] = struct{}{}
	}

	for worldID := range uniqueWorldIDs {
		if err := openWorldPage(worldID); err != nil {
			fmt.Fprintf(os.Stderr, "URLを開けませんでした: %s: %v\n", worldID, err)
			hadError = true
			continue
		}
		openedCount++
	}

	if openedCount == 0 {
		fmt.Println("有効なWorldID情報が見つかりませんでした.")
		hadError = true
	}

	// 途中でエラーがあった場合、メッセージを確認用に待機
	if hadError {
		waitForEnter()
	}
}
