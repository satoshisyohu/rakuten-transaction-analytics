package impl

import (
	"bytes"
	"encoding/csv"
	"mime/multipart"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

type FileReaderRule struct {
}

func NewFileReaderRule() *FileReaderRule {
	return &FileReaderRule{}
}

func (f *FileReaderRule) ReadAll(fileHeader *multipart.FileHeader) ([][]string, error) {

	file, err := fileHeader.Open()
	// // エラーチェック
	if err != nil {
		return nil, err
	}
	defer file.Close()
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
	_, err := file.Read(bom)
	if err != nil {
		return nil, err
	}

	// bomの判定ロジックを実施する
	expectedBOM := []byte{0xEF, 0xBB, 0xBF}
	// UTF-8のBOMか判定する
	if bytes.Equal(bom, expectedBOM) {
		// BOMが存在する場合はすでにファイルポインタがBOM以降のためそのまま返す
		// これはクレジットカードの明細を読み込んだ場合でありUTF-8のためそのまま読み込み可能
		return csv.NewReader(file), nil
	} else {
		// BOMがない場合、ファイルポインタを先頭に戻す必要がある
		// これはデビットカードの明細を読み込んだ場合でありShift-JISのためデコードして読み込む必要がある
		_, err = file.Seek(0, 0)
		if err != nil {
			return nil, err
		}
		// BOMがない場合はShift-JISでデコードする
		return csv.NewReader(transform.NewReader(file, japanese.ShiftJIS.NewDecoder())), nil
	}
}
