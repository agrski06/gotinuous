package internal

import (
	"bytes"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	Variables  map[string]string `yaml:"variables"`
	Stages     yaml.MapSlice     `yaml:"stages"`
	WorkingDir string
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

	wd, err := os.Getwd()
	if err != nil {
		panic("Could not get current working directory")
	}
	tool.WorkingDir = wd

	return tool
}

func (tool Tool) InitRepository() {
	if tool.Conf.Repository.URL == "" {
		log.Println("No repository specified. Skipping repository initialization.")
		return
	}

	repositoryNameCommand := exec.Command("basename", tool.Conf.Repository.URL)
	repositoryName := strings.Split(handleCommand(repositoryNameCommand), ".")[0]

	log.Println("Cloning", repositoryName, "repository at", tool.Conf.Repository.URL)

	newWorkingDirectory := filepath.Join(tool.WorkingDir, repositoryName)
	exists, _ := pathExists(newWorkingDirectory)
	if exists {
		log.Println("Repository", repositoryName, "already exists, skipping...")
		return
	}
	tool.WorkingDir = newWorkingDirectory

	gitCommand := exec.Command("git", "clone", tool.Conf.Repository.URL)
	handleCommand(gitCommand)

	log.Println("Cloned git repository")
}

func (tool Tool) LoadVariablesIntoEnv(env []string) {
	log.Println("Loading variables to env")
	for k, v := range tool.Variables {
		env = append(env, k+"="+v)
	}
}

func (tool Tool) ExecStages(env []string) {
	for _, stageMap := range tool.Stages {
		stage := convertMapSliceToStage(stageMap.Value.(yaml.MapSlice))
		wd := filepath.Join(tool.WorkingDir, stage.Dir)
		log.Print("Executing stage: ", stageMap.Key.(string), " (", wd, ")")

		// replace variables in command
		// exec command
	}
}

func convertMapSliceToStage(slice yaml.MapSlice) Stage {
	var stage Stage

	for _, item := range slice {
		key := item.Key.(string)
		value := item.Value

		switch key {
		case "dir":
			stage.Dir = value.(string)
		case "command":
			stage.Command = value.(string)
		}
	}

	return stage
}

func handleCommand(command *exec.Cmd) string {
	var stderr bytes.Buffer
	command.Stderr = &stderr
	output, err := command.Output()
	if err != nil {
		panic(stderr.String())
	}
	return string(output)
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
