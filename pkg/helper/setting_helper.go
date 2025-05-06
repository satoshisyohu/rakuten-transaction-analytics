package helper

import (
	"fmt"
)

// todo 今後は環境変数から読み込むようにする
var _ = "SETTING_PATH"
var label = &Setting{}

// Load 設定を読み込む
func Load() error {
	// todo pathは環境変数から読み込む
	err := JsonToStruct("../pkg/resource/setting.yaml", label)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// Setting yamlから読み込んだ設定にmappingする構造体
type Setting struct {
	Bigquery Bigquery `yaml:"bigquery"`
}

// Bigquery bqの設定に関する構造体
type Bigquery struct {
	ProjectId string `yaml:"projectId"`
	DatasetId string `yaml:"datasetId"`
}

// GetProjectId ProjectIdを取得する

func GetProjectId() string {
	return label.Bigquery.ProjectId
}

// GetDatasetId DatasetIdを取得する

func GetDatasetId() string {
	return label.Bigquery.DatasetId
}
