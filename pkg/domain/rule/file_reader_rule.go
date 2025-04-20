package rule

import "mime/multipart"

type IFileReaderRule interface {
	// ReadAll ファイルを読み込む
	ReadAll(file *multipart.FileHeader) ([][]string, error)
}
