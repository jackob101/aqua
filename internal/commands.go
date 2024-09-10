package internal

import (
	"encoding/json"
	"errors"
	"io/fs"
	"jackob101/aqua/common/dto"
	"log/slog"
	"os"
	"path"
)

// TODO: This also should be read from env variables
var (
	commandFilename     = "lake.json"
	howFarUpShouldCheck = 3
)

// TODO: add error handling to this function. Should also probably add error pane to nicely display these errors to user
func ReadCommands() ([]dto.Command, error) {
	dir, err := os.Getwd()
	if err != nil {
		return []dto.Command{}, nil
	}

	foundGitFiles := false

	for i := 0; i < howFarUpShouldCheck; i++ {
		contains, err := containsGitFiles(dir)
		if err != nil {
			return nil, err
		}
		if !contains {
			dir, _ = path.Split(dir)
			continue
		}
		foundGitFiles = true
	}

	if !foundGitFiles {
		return nil, errors.New("couldn't find git files")
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, errors.New("failed to read directory: " + err.Error())
	}

	var commandFile fs.DirEntry

	for _, e := range entries {
		if e.Name() == commandFilename {
			commandFile = e
		}
	}
	if commandFile == nil {
		return nil, errors.New("lake.json is missing")
	}
	var content []byte
	content, err = os.ReadFile(path.Join(dir, commandFile.Name()))
	if err != nil {
		return nil, errors.New("failed to read command file: " + err.Error())
	}

	var commands []dto.Command
	err = json.Unmarshal(content, &commands)
	if err != nil {
		return nil, errors.New("failed to parse command file: " + err.Error())
	}
	return commands, nil
}

func containsGitFiles(path string) (bool, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		slog.Error("Failed to read directory", "error", err.Error())
		return false, errors.New("failed to read directory: " + err.Error())
	}

	for _, e := range entries {
		if e.Name() == ".git" {
			return true, nil
		}
	}
	return false, nil
}
