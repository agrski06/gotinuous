package internal

import (
	"bytes"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type Stage struct {
	Dir     string `yaml:"dir"`
	Command string `yaml:"command"`
}

type Tool struct {
	Conf struct {
		Repository struct {
			URL string `yaml:"url"`
		} `yaml:"repository"`
	} `yaml:"conf"`
	Variables map[string]string `yaml:"variables"`
	Stages    map[string]Stage  `yaml:"stages"`
}

func InitTool() Tool {
	filename, err := filepath.Abs("./config.yaml")
	if err != nil {
		panic(err)
	}

	file, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	var tool Tool
	err = yaml.Unmarshal(file, &tool)
	if err != nil {
		panic(err)
	}

	return tool
}

func (tool Tool) InitRepository() {
	log.Println("Cloning git repository", tool.Conf.Repository.URL)

	gitCommand := exec.Command("git", "clone", tool.Conf.Repository.URL)
	var stderr bytes.Buffer
	gitCommand.Stderr = &stderr
	_, err := gitCommand.Output()
	if err != nil {
		panic(stderr.String())
	}
	log.Println("Cloned git repository")
}
