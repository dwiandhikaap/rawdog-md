package helper

import (
	"os"
	"path/filepath"
)

func TruncateString(s string, n int) string {
	i := 0
	for j := range s {
		if i == n {
			return s[:j] + "..."
		}
		i++
	}
	return s
}

func OmitFilenameExtension(filename string) string {
	return filename[:len(filename)-len(filepath.Ext(filename))]
}

func WriteHtmlFile(absolutePath string, content string) error {
	dir := filepath.Dir(absolutePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	if _, err := os.Stat(absolutePath); err == nil {
		err := os.Remove(absolutePath)
		if err != nil {
			return err
		}
	}

	file, err := os.Create(absolutePath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}
