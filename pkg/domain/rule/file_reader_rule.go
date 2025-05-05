//go:generate mockgen -source=file_reader_rule.go -destination=../../../internal/mocks/domain/rule/file_reader_rule.go -package=mocks

package rule

import "mime/multipart"

type IFileReaderRule interface {
	// ReadAll ファイルを読み込む
	ReadAll(file *multipart.FileHeader) ([][]string, error)
}
