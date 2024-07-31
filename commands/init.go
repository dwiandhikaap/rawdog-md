package commands

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"rawdog-md/presets"
)

func Init(relativePath string, preset string) error {
	err := os.Mkdir(relativePath, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.Chdir(relativePath)
	if err != nil {
		return err
	}

	var presetFs embed.FS
	switch preset {
	case "basic":
		presetFs = presets.BasicPreset
	default:
		presetFs = presets.SkeletonPreset
	}

	err = fs.WalkDir(presetFs, preset, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Ignore build folder
		buildDir := filepath.Join(preset, "build")
		if filepath.Clean(path) == buildDir {
			return filepath.SkipDir
		}

		if d.IsDir() {
			// Make directory relative to preset
			relativePath, err := filepath.Rel(preset, path)
			if err != nil {
				return err
			}
			if relativePath == "." {
				return nil
			}
			return os.Mkdir(relativePath, os.ModePerm)
		}

		// Open file
		file, err := presetFs.Open(path)
		if err != nil {
			return err
		}

		fileInfo, err := file.Stat()
		if err != nil {
			return err
		}

		fileContent := make([]byte, fileInfo.Size())
		_, err = file.Read(fileContent)
		if err != nil {
			return err
		}

		// Get relative target path
		targetPath, err := filepath.Rel(preset, path)
		if err != nil {
			return err
		}

		// Get absolute target path from current working directory
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		abs := filepath.Join(cwd, targetPath)

		// Create file
		err = os.MkdirAll(filepath.Dir(abs), os.ModePerm)
		if err != nil {
			return err
		}

		// Write file
		err = os.WriteFile(abs, fileContent, os.ModePerm)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}