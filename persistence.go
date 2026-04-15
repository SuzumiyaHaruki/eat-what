package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

const appID = "com.example.eatwhat"

type appState struct {
	Options     []string `json:"options"`
	MenuName    string   `json:"menu_name"`
	MenuPath    string   `json:"menu_path"`
	CurrentFile string   `json:"current_file"`
}

func loadAppState() (appState, error) {
	path, err := appStatePath()
	if err != nil {
		return appState{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return appState{}, nil
		}
		return appState{}, err
	}

	var state appState
	if err := json.Unmarshal(data, &state); err != nil {
		return appState{}, err
	}
	return state, nil
}

func saveAppState(state appState) error {
	path, err := appStatePath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}

func appStatePath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "eat-what", "state.json"), nil
}

func menusDirPath() (string, error) {
	workdir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return workdir, nil
}

func managedMenuPath(name string) (string, error) {
	dir, err := menusDirPath()
	if err != nil {
		return "", err
	}
	fileName := sanitizeMenuFileName(name)
	if fileName == "" {
		fileName = "untitled-menu"
	}
	return filepath.Join(dir, fileName+".txt"), nil
}

func saveOptionsToTxt(path string, options []string) error {
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	content := strings.Join(options, "\n")
	if content != "" {
		content += "\n"
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

func ensureMenuFile(path string) error {
	if path == "" {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	return file.Close()
}

func renameManagedMenuFile(oldPath, newPath string, options []string) (string, error) {
	if newPath == "" {
		return "", nil
	}
	if oldPath == "" || oldPath == newPath {
		return newPath, saveOptionsToTxt(newPath, options)
	}
	if err := os.MkdirAll(filepath.Dir(newPath), 0o755); err != nil {
		return "", err
	}
	if err := os.Rename(oldPath, newPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return "", err
	}
	if err := saveOptionsToTxt(newPath, options); err != nil {
		return "", err
	}
	return newPath, nil
}

func sanitizeMenuFileName(name string) string {
	name = strings.TrimSpace(name)
	replacer := strings.NewReplacer(
		"/", "-",
		"\\", "-",
		":", "-",
		"*", "-",
		"?", "",
		"\"", "",
		"<", "",
		">", "",
		"|", "-",
	)
	name = replacer.Replace(name)
	name = strings.Trim(name, ". ")
	return name
}
