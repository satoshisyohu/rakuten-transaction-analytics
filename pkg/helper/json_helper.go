package helper

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

func JsonToStruct[T any](path string, v *T) (err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}
	if err = yaml.Unmarshal(data, v); err != nil {
		return fmt.Errorf("unmarshal err: %w", err)
	}
	return nil
}
