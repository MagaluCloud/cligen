package gen_cli_code

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const FILE_CUSTOM_COMMANDS = "base-cli-custom/custom-commands.yaml"

type (
	CustomData struct {
		FileID string `yaml:"file_id"`
		Use    string `yaml:"use"`
	}

	CustomHeader struct {
		Data []CustomData `yaml:"data"`
	}
)

func NewCustom() *CustomHeader {
	return &CustomHeader{Data: []CustomData{}}
}

func (c *CustomHeader) Add(custom CustomData) error {
	c.Data = append(c.Data, custom)
	return nil
}

func (c *CustomHeader) Write() error {
	content, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(FILE_CUSTOM_COMMANDS, content, 0644)
}

func (c *CustomHeader) Load() error {
	content, err := os.ReadFile(FILE_CUSTOM_COMMANDS)
	if err != nil {
		fmt.Println("error a")
		fmt.Println(FILE_CUSTOM_COMMANDS)
		return err
	}
	var customUses CustomHeader
	err = yaml.Unmarshal(content, &customUses)
	if err != nil {
		fmt.Println("error b")
		fmt.Println(FILE_CUSTOM_COMMANDS)
		return err
	}
	c.Data = customUses.Data
	return nil
}

func (c *CustomHeader) Find(fileID string) *CustomData {
	for _, custom := range c.Data {
		if custom.FileID == fileID {
			return &CustomData{FileID: custom.FileID, Use: custom.Use}
		}
	}
	return nil
}

func NewCustomData(fileID string) *CustomData {
	return &CustomData{FileID: fileID}
}

func (c *CustomData) AddUse(use string) {
	c.Use = use
}
