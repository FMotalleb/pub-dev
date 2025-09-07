package pub

import (
	"errors"

	"gopkg.in/yaml.v3"

	"github.com/fmotalleb/pub-dev/tar"
)

type Spec struct {
	Name    string
	Version string
	Raw     map[string]any
}

func ReadSpecFromTar(file string) (*Spec, error) {
	data, err := tar.ReadFile(file, "pubspec.yaml")
	if err != nil {
		return nil, err
	}
	return ParseSpec(data)
}

func ParseSpec(content []byte) (*Spec, error) {
	raw := make(map[string]any, 0)
	err := yaml.Unmarshal(content, raw)
	spec := new(Spec)
	var ok bool
	if spec.Name, ok = getString(raw, "name"); !ok {
		return nil, errors.New("invalid pubspec data: name")
	}
	if spec.Version, ok = getString(raw, "version"); !ok {
		return nil, errors.New("invalid pubspec data: version")
	}
	spec.Raw = raw
	return spec, err
}

func getString(src map[string]any, key string) (string, bool) {
	var value any
	var ok bool
	if value, ok = src[key]; !ok {
		return "", false
	}
	var str string
	if str, ok = value.(string); !ok {
		return "", false
	}
	return str, true
}
