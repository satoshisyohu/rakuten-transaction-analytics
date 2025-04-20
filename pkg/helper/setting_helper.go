package helper

import (
	"fmt"
)

// todo 今後は環境変数から読み込むようにする
var _ = "SETTING_PATH"
var label = &Setting{}

func Load() error {
	// todo pathは環境変数から読み込む
	err := JsonToStruct("../pkg/resource/setting.yaml", label)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

type Setting struct {
	Bigquery Bigquery `yaml:"bigquery"`
}

type Bigquery struct {
	ProjectId string `yaml:"projectId"`
	DatasetId string `yaml:"datasetId"`
}

func GetProjectId() string {
	return label.Bigquery.ProjectId
}

func GetDatasetId() string {
	return label.Bigquery.DatasetId
}
