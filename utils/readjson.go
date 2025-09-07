package utils

import (
	"encoding/json"
	"os"

	"github.com/fmotalleb/go-tools/template"
)

const filePermission = 0o600

func ReadJSON(path string, dst any) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(content, dst)
}

func ReadJSONTemplate(path string, dst any, data ...any) error {
	var extra any
	if len(data) > 0 {
		extra = data[0]
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	strContent, err := template.EvaluateTemplate(string(content), extra)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(strContent), dst)
}

func WriteJSON(path string, src any) error {
	content, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return os.WriteFile(path, content, filePermission)
}
