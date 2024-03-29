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

type Repository struct {
	URL string `yaml:"url"`
}

type Tool struct {
	Conf struct {
		Repository Repository `yaml:"repository"`
	} `yaml:"conf"`
	Variables map[string]string `yaml:"variables"`
	Stages    yaml.MapSlice     `yaml:"stages"`
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
	if tool.Conf.Repository.URL == "" {
		log.Println("No repository specified. Skipping repository initialization.")
		return
	}

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

func (tool Tool) LoadVariablesIntoEnv(env []string) {
	for k, v := range tool.Variables {
		env = append(env, k+"="+v)
	}
}

func (tool Tool) ExecStages(env []string) {
	for _, stage := range tool.Stages {
		log.Print("Executing stage: " + stage.Key.(string))

	}
}
