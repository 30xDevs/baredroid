package baredroid

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type PkgInstall struct {
	Name	string `json:"name"`
	Package string `json:"package"`
	Type string `json:"type"`
	Source 	string `json:"source"`
	Children []PkgInstall `json:"children"`
}

type Config struct {
	PkgRemove []string `json:"pkgRemove"`
	PkgInstall []PkgInstall `json:"pkgInstall"`
}

func NewConfig(ConfPath string) (*Config, error) {

	configFile, err := os.Open(ConfPath)
	
	if err != nil {
		return nil, fmt.Errorf("could not open config at %s: %s", ConfPath, err)
	}

	defer configFile.Close()

	byteValue, _ := io.ReadAll(configFile)

	var config Config
	json.Unmarshal(byteValue, &config)

	return &config, nil
}