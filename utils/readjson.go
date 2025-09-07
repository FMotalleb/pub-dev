package utils

import (
	"encoding/json"
	"os"

	"github.com/fmotalleb/go-tools/template"
)

func ReadJSON(path string, dst any) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(content, dst); err != nil {
		return err
	}
	return nil
}

func ReadJSONTemplate(path string, dst any, data ...any) error {
	var extra any
	if len(data) != 0 {
		extra = data[0]
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	strContent, err := template.EvaluateTemplate(
		string(content),
		extra,
	)
	if err != nil {
		return err
	}
	if err = json.Unmarshal([]byte(strContent), dst); err != nil {
		return err
	}
	return nil
}
