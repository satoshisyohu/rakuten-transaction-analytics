package impl

import (
	"bytes"
	"encoding/csv"
	"log/slog"
	"mime/multipart"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// FileReaderRule CSVファイルを読み込むルール
type FileReaderRule struct {
}

// NewFileReaderRule FileReaderRuleのファクトリ関数
func NewFileReaderRule() *FileReaderRule {
	return &FileReaderRule{}
}

// ReadAll ファイルを読み込む
func (f *FileReaderRule) ReadAll(fileHeader *multipart.FileHeader) ([][]string, error) {

	file, err := fileHeader.Open()
	// // エラーチェック
	if err != nil {
		return nil, err
	}
	defer func(file multipart.File) {
		dErr := file.Close()
		if dErr != nil {
			slog.Error("ファイルクローズに失敗しました", "err", dErr)
		}
	}(file)
	// CSVリーダーを作成
	reader, err := f.newReader(file)
	if err != nil {
		return nil, err
	}
	// CSVを読み込む
	return reader.ReadAll()
}

func (f *FileReaderRule) newReader(file multipart.File) (*csv.Reader, error) {

	// BOM判定のためファイルの先頭から3バイトを読み込む
	bom := make([]byte, 3)
	if _, err := file.Read(bom); err != nil {
		return nil, err
	}

	// bomの判定ロジックを実施する
	// UTF-8のBOMか判定する
	if bytes.Equal(bom, []byte{0xEF, 0xBB, 0xBF}) {
		// BOMが存在する場合はすでにファイルポインタがBOM以降のためそのまま返す
		// これはクレジットカードの明細を読み込んだ場合でありUTF-8のためそのまま読み込み可能
		return csv.NewReader(file), nil
	}
	// BOMがない場合、ファイルポインタを先頭に戻す必要がある
	// これはデビットカードの明細を読み込んだ場合でありShift-JISのためデコードして読み込む必要がある
	if _, err := file.Seek(0, 0); err != nil {
		return nil, err
	}
	// BOMがない場合はShift-JISでデコードする
	return csv.NewReader(transform.NewReader(file, japanese.ShiftJIS.NewDecoder())), nil
	
}
