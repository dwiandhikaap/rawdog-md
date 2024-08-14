package helper

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	"github.com/tdewolff/minify/json"
	"github.com/tdewolff/minify/svg"
	"github.com/tdewolff/minify/xml"
)

var Minifier = minify.New()

func init() {
	Minifier.AddFunc("text/css", css.Minify)
	Minifier.AddFunc("text/html", html.Minify)
	Minifier.AddFunc("image/svg+xml", svg.Minify)
	Minifier.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	Minifier.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	Minifier.AddFuncRegexp(regexp.MustCompile("[/+]xml$"), xml.Minify)
}

func SliceContainsInt(slice []int, n int) bool {
	for _, v := range slice {
		if v == n {
			return true
		}
	}
	return false
}

func SliceContainsString(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

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

func WriteTextFile(absolutePath string, content string) error {
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

func IsPathDir(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fi.IsDir()
}
