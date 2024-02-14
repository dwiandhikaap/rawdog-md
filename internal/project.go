package internal

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"rawdog-md/global"
	"rawdog-md/helper"
)

func LoadProject() (*[]Page, error) {
	pagesDir, err := LoadPages()
	if err != nil {
		return nil, err
	}

	templates := make(map[string]*Template)
	pages := make([]Page, 0)

	for _, dir := range *pagesDir {
		page, err := NewPage(dir)
		if err != nil {
			return nil, fmt.Errorf("error parsing page '%s': %v", dir, err)
		}
		if page == nil {
			continue
		}

		if page.TemplateName != nil {
			if _, ok := templates[*page.TemplateName]; !ok {
				template, err := LoadTemplate(global.Config.RootRelativePath, *page.TemplateName)
				if err != nil {
					return nil, fmt.Errorf("error loading template '%s': %v", *page.TemplateName, err)
				}
				templates[*page.TemplateName] = template
			}

			page.Template = templates[*page.TemplateName]
		}

		pages = append(pages, *page)
	}

	contexts := NewContexts(pages)

	// Render
	for i := range pages {
		page := &pages[i]
		err := page.Render(contexts)
		if err != nil {
			return nil, fmt.Errorf("error rendering page '%s': %v", page.SourceAbsolutePath, err)
		}
		fmt.Println(&page.SourceAbsolutePath)
	}

	return &pages, nil
}

func LoadPages() (*[]string, error) {
	pages := make([]string, 0)

	pagesDirAbs := filepath.Join(global.Config.RootAbsolutePath, "pages")

	err := filepath.WalkDir(pagesDirAbs, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		pages = append(pages, path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &pages, nil
}

func LoadTemplate(rootRelativePath string, templateName string) (*Template, error) {
	templatePath := filepath.Join(rootRelativePath, "templates", templateName)
	templatePathAbs, err := filepath.Abs(templatePath)
	if err != nil {
		return nil, err
	}

	fileContent, err := os.ReadFile(templatePathAbs)
	if err != nil {
		return nil, err
	}

	return &Template{
		AbsolutePath: templatePathAbs,
		Filename:     templateName,
		Content:      string(fileContent),
	}, nil
}

func WritePages(pages *[]Page) error {
	for _, page := range *pages {
		abs := filepath.Join(global.Config.RootAbsolutePath, "build", page.RelativeUrl)
		err := helper.WriteHtmlFile(abs, page.Output)
		if err != nil {
			return err
		}
	}

	return nil
}

func CopyStaticFiles() error {
	staticDirAbs := filepath.Join(global.Config.RootAbsolutePath, "static")
	buildDirAbs := filepath.Join(global.Config.RootAbsolutePath, "build")

	err := filepath.WalkDir(staticDirAbs, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(staticDirAbs, path)
		if err != nil {
			return err
		}

		abs := filepath.Join(buildDirAbs, rel)
		err = os.MkdirAll(filepath.Dir(abs), os.ModePerm)
		if err != nil {
			return err
		}

		fileContent, err := os.ReadFile(path)
		if err != nil {
			return err
		}

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
