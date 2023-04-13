package main

// 実行したら指定ディレクトリにファイルを作成する
import (
	_ "embed"
	"flag"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

const (
	// config_dir_name
	configName = "go-memo"
)

var (
	filename string // ファイル名

	//go:embed template.md
	tmpl string
)

type Memo struct {
	// ファイル名
	Filename string
	// 日付
	Date string
}

func init() {
	flag.StringVar(&filename, "n", "default.md", "failenameを指定")
}

func main() {
	// フラグ解析
	flag.Parse()

	// XDG_CONFIG_HOMEが設定されているか確認する
	// 設定されていない場合は、ホームディレクトリに作成する
	xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfigHome == "" {
		// ホームディレクトリのパスを取得する
		home := getHomeDir()
		xdgConfigHome = filepath.Join(home, ".config")
	}

	// ディレクトリを作成する
	// ディレクトリ名は、XDG_CONFIG_HOME/go-memo
	c_dir := filepath.Join(xdgConfigHome, configName)
	if _, err := os.Stat(c_dir); os.IsNotExist(err) {
		os.MkdirAll(c_dir, 0755)
	}

	// ファイルを作成する
	// ファイル名は、引数で指定されたもの
	// ファイルの中身は、テンプレート
	memo := Memo{
		Filename: filename,
		Date:     getTime(),
	}

	// ファイルを作成する
	// xdgConfigHome/go-memo/template.mdがあれば、それを使う
	// なければ、埋め込みのテンプレートを使う
	var t *template.Template
	if _, err := os.Stat(filepath.Join(c_dir, "template.md")); os.IsExist(err) {

		t = template.Must(template.ParseFiles(filepath.Join(c_dir, "template.md")))

	} else {

		t, err = template.New("sample").Parse(tmpl)
		if err != nil {
			log.Fatal(err)
		}
	}

	// todo: ファイル名が重複していたら、エラーを出す
	// todo: ファイルの出力先をXDG_CONFIG_HOME/go-memoにする
	if err := t.Execute(os.Stdout, memo); err != nil {
		log.Fatal(err)
	}

}

func getTime() string {
	// 時刻を取得する
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
