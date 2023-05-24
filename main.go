package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

const (
	configDirName = "go-memo"
)

var (
	//go:embed template.md
	tmpl string
)

type Memo struct {
	Filename string
	Date     string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("ファイル名を指定してください")
		os.Exit(1)
	}

	filename := os.Args[1]

	if err := run(filename); err != nil {
		fmt.Printf("エラーが発生しました: %v\n", err)
		os.Exit(1)
	}
}

func run(filename string) error {
	err := validateFilename(filename)
	if err != nil {
		return fmt.Errorf("ファイル名の検証に失敗しました: %w", err)
	}

	err = createMemoFile(filename)
	if err != nil {
		return fmt.Errorf("メモファイルの作成に失敗しました: %w", err)
	}

	err = openFileInEditor(filename)
	if err != nil {
		return fmt.Errorf("エディタでファイルを開く際にエラーが発生しました: %w", err)
	}

	return nil
}

func validateFilename(filename string) error {
	if err := validateEmpty(filename); err != nil {
		return err
	}

	if err := validateLength(filename); err != nil {
		return err
	}

	if err := validateInvalidChars(filename); err != nil {
		return err
	}

	return nil
}

func validateEmpty(filename string) error {
	// ファイル名が空文字列でないことを確認
	if strings.TrimSpace(filename) == "" {
		return fmt.Errorf("ファイル名が空です")
	}
	return nil
}

func validateLength(filename string) error {
	// ファイル名が30文字以下であることを確認
	if len(filename) > 30 {
		return fmt.Errorf("ファイル名が30文字を超えています")
	}
	return nil
}

func validateInvalidChars(filename string) error {
	// ファイル名に使用できない文字が含まれていないことを確認
	invalidChars := []string{`\`, `/`, `:`, `*`, `?`, `"`, `<`, `>`, `|`}
	for _, char := range invalidChars {
		if strings.Contains(filename, char) {
			return fmt.Errorf("ファイル名に使用できない文字が含まれています: %s", char)
		}
	}
	return nil
}

func createMemoFile(filename string) error {
	// Configuration directory path
	configDir, err := getOrCreateConfigDir()
	if err != nil {
		return err
	}

	// Prepare template
	t, err := prepareTemplate(configDir)
	if err != nil {
		return err
	}

	// Memo directory path
	memoDir, err := getOrCreateMemoDir(configDir)
	if err != nil {
		return err
	}

	// Full path of the file to be created
	filenameFull := filepath.Join(memoDir, filename)

	// Check if file already exists
	if _, err := os.Stat(filenameFull); !os.IsNotExist(err) {
		if err == nil {
			return fmt.Errorf("file already exists: %s", filename)
		} else {
			return fmt.Errorf("failed to get file status: %v", err)
		}
	}

	// Create and write to file
	f, err := os.Create(filenameFull)
	if err != nil {
		return fmt.Errorf("failed to create memo file: %v", err)
	}
	defer f.Close()

	memo := Memo{
		Filename: filename,
		Date:     getTime(),
	}

	if err := t.Execute(f, memo); err != nil {
		return fmt.Errorf("failed to write to memo file: %v", err)
	}

	fmt.Println("File creation completed:", filenameFull)
	return nil
}

func getOrCreateConfigDir() (string, error) {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		home := getHomeDir()
		xdgConfigHome = filepath.Join(home, ".config")
	}
	configDir := filepath.Join(xdgConfigHome, configDirName)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err = os.MkdirAll(configDir, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create config directory: %v", err)
		}
	}
	return configDir, nil
}

func prepareTemplate(configDir string) (*template.Template, error) {
	var t *template.Template
	if _, err := os.Stat(filepath.Join(configDir, "template.md")); os.IsExist(err) {
		t, err = template.ParseFiles(filepath.Join(configDir, "template.md"))
		if err != nil {
			return nil, fmt.Errorf("failed to load template: %v", err)
		}
	} else {
		t, err = template.New("sample").Parse(tmpl)
		if err != nil {
			return nil, fmt.Errorf("failed to parse template: %v", err)
		}
	}
	return t, nil
}

func getOrCreateMemoDir(configDir string) (string, error) {
	memoDir := filepath.Join(configDir, "memo")
	if _, err := os.Stat(memoDir); os.IsNotExist(err) {
		err = os.MkdirAll(memoDir, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create memo directory: %v", err)
		}
	}
	return memoDir, nil
}

func getTime() string {
	// 現在の時刻を取得する
	t := time.Now()
	return t.Format("2006-01-02-15:04:05")
}

// ホームディレクトリのパスを取得する
func getHomeDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	return home
}

func openFileInEditor(filename string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return fmt.Errorf("環境変数EDITORが設定されていません")
	}

	filenameFull := getMemoFilePath(filename)
	cmd := exec.Command(editor, filenameFull)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("エディタでファイルを開く際にエラーが発生しました: %v", err)
	}

	return nil
}

func getMemoFilePath(filename string) string {
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		home := getHomeDir()
		xdgConfigHome = filepath.Join(home, ".config")
	}
	configDir := filepath.Join(xdgConfigHome, configDirName)
	memoDir := filepath.Join(configDir, "memo")
	return filepath.Join(memoDir, filename)
}
