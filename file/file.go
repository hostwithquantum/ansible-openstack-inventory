package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// GroupVarsFile ...
type GroupVarsFile struct {
	Path string
}

// NewGroupVarsFile ... factory/ctor
func NewGroupVarsFile(path string) GroupVarsFile {
	f := GroupVarsFile{
		Path: path,
	}
	return f
}

// HandleGroup ...
func (gvf GroupVarsFile) HandleGroup(group string) (map[string]interface{}, error) {
	groupFile := fmt.Sprintf("%s.yml", group)
	fileName := filepath.Join(gvf.Path, groupFile)

	info, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, err
	}

	if info.IsDir() {
		return nil, fmt.Errorf("Path is a directory: %s", fileName)
	}

	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	groupFileYaml := make(map[string]interface{})
	err = yaml.Unmarshal(buf, &groupFileYaml)
	if err != nil {
		return nil, err
	}

	return groupFileYaml, nil
}
