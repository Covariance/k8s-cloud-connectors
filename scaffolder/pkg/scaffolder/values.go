// Copyright (c) 2021 Yandex LLC. All rights reserved.
// Author: Martynov Pavel <covariance@yandex-team.ru>

package scaffolder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v3"
)

// Values is a container with strings that are to be substituted in scaffolding and scheme.
type Values map[string]string

func ParseValuesByUnmarshal(path string, unmarshal func([]byte, interface{}) error) (Values, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open value file: %w", err)
	}

	var res map[string]string

	if err := unmarshal(contents, &res); err != nil {
		return nil, fmt.Errorf("unable to parse value file: %w", err)
	}
	return res, nil
}

func ParseValuesFromJSON(path string) (Values, error) {
	return ParseValuesByUnmarshal(path, json.Unmarshal)
}

func ParseValuesFromYAML(path string) (Values, error) {
	return ParseValuesByUnmarshal(path, yaml.Unmarshal)
}

func ParseValuesFromFile(path string) (Values, error) {
	if strings.HasSuffix(path, ".json") {
		return ParseValuesFromJSON(path)
	}
	if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		return ParseValuesFromYAML(path)
	}

	return nil, fmt.Errorf("unable to detect type of file \"%s\"", path)
}

func ParseValuesFromString(contents string) (Values, error) {
	split := strings.SplitN(contents, "=", 2)
	if len(split) == 1 {
		if strings.HasPrefix(contents, "=") {
			return nil, fmt.Errorf("invalid \"key=value\" pair provided: %s", contents)
		}
		return map[string]string{split[0]: ""}, nil
	}
	return map[string]string{split[0]: split[1]}, nil
}
