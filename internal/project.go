package internal

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/dwiandhikaap/rawdog-md/global"
	"github.com/dwiandhikaap/rawdog-md/helper"
)

type Project struct {
	Pages  []Page
	Assets []string
}

func (p *Project) WritePages() error {
	err := WritePages(&p.Pages)
	if err != nil {
		return err
	}

	return nil
}

func (p *Project) CopyStaticFiles() error {
	err := CopyStaticFiles()
	if err != nil {
		return err
	}

	return nil
}

func (p *Project) PurgeBuildDir() error {
	buildDirAbs := filepath.Join(global.Config.RootAbsolutePath, "build")

	// Check if build directory exists
	_, err := os.Stat(buildDirAbs)
	if os.IsNotExist(err) {
		return nil
	}

	// Check if dir is not important to prevent accidental deletion
	if buildDirAbs == "/" || buildDirAbs == "~" || buildDirAbs == "" {
		return fmt.Errorf("build directory is important, refusing to delete")
	}

	// More check so we dont fuck up like valve did.. hopefully...
	actualRootAbs, err := filepath.Abs(global.Config.RootRelativePath)
	if err != nil {
		return err
	}

	actualRootAbs = strings.ReplaceAll(actualRootAbs, "\\", "/")

	if actualRootAbs != global.Config.RootAbsolutePath {
		return fmt.Errorf("project root path \"%s\", doesnt not equal current working dir \"%s\"", global.Config.RootAbsolutePath, actualRootAbs)
	}

	// scaryyyyyy
	err = os.RemoveAll(buildDirAbs)
	if err != nil {
		return err
	}

	return nil
}

func (p *Project) ForceRebuild() error {
	// O God, I humbly implore Thee to forgive the cardinal sin I am about to commit against the garbage collector and disk I/O,
	// as I shamefully recreating the entire project anew with each file change, without reusing objects or memory.
	pages, err := LoadPages()
	if err != nil {
		return err
	}

	p.Pages = *pages

	err = p.PurgeBuildDir()
	if err != nil {
		return err
	}

	p.WritePages()
	p.CopyStaticFiles()

	return nil
}

func NewProject() (*Project, error) {
	return &Project{}, nil
}

func LoadPages() (*[]Page, error) {
	pagesDir, err := GetPagesPath()
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
				template, err := LoadTemplate(global.Config.RootAbsolutePath, *page.TemplateName)
				if err != nil {
					return nil, fmt.Errorf("error loading template '%s': %v", *page.TemplateName, err)
				}
				templates[*page.TemplateName] = template
			}

			page.Template = templates[*page.TemplateName]
		}

		pages = append(pages, *page)
	}

	// Check any conflicting page paths
	// {"path": ["filename1", "filename2", "filename3", ...]}
	paths := make(map[string][]string)
	for _, page := range pages {
		if _, ok := paths[page.RelativeUrl]; !ok {
			paths[page.RelativeUrl] = make([]string, 0)
		}
		paths[page.RelativeUrl] = append(paths[page.RelativeUrl], page.SourceAbsolutePath)
	}

	for path, files := range paths {
		if len(files) > 1 {
			sourcePaths := ""
			for _, file := range files {
				sourcePaths += fmt.Sprintf("\n	'%s', ", file)
			}

			return nil, fmt.Errorf("conflicting page paths: the following files have the same output path '%s': %s", path, sourcePaths)
		}
	}

	contexts := NewContexts(pages)

	// Render
	for i := range pages {
		page := &pages[i]
		err := page.Render(contexts)
		if err != nil {
			return nil, fmt.Errorf("error rendering page '%s': %v", page.SourceAbsolutePath, err)
		}
	}

	return &pages, nil
}

func GetPagesPath() (*[]string, error) {
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

func LoadTemplate(root string, templateName string) (*Template, error) {
	templatePathAbs := filepath.Join(root, "templates", templateName)

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
		err := helper.WriteTextFile(abs, page.Output)
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

		// TODO: Add more options for each minification
		fileExtension := filepath.Ext(path)
		if fileExtension == ".html" && global.Config.UserConfig.Options.Minify.HTML {
			fileContent, err = helper.Minifier.Bytes("text/html", fileContent)
			if err != nil {
				return err
			}
		} else if fileExtension == ".css" && global.Config.UserConfig.Options.Minify.CSS {
			fileContent, err = helper.Minifier.Bytes("text/css", fileContent)
			if err != nil {
				return err
			}
		} else if fileExtension == ".js" && global.Config.UserConfig.Options.Minify.JS {
			fileContent, err = helper.Minifier.Bytes("application/javascript", fileContent)
			if err != nil {
				return err
			}
		} else if fileExtension == ".json" && global.Config.UserConfig.Options.Minify.JSON {
			fileContent, err = helper.Minifier.Bytes("application/json", fileContent)
			if err != nil {
				return err
			}
		} else if fileExtension == ".xml" && global.Config.UserConfig.Options.Minify.XML {
			fileContent, err = helper.Minifier.Bytes("text/xml", fileContent)
			if err != nil {
				return err
			}
		} else if fileExtension == ".svg" && global.Config.UserConfig.Options.Minify.SVG {
			fileContent, err = helper.Minifier.Bytes("image/svg+xml", fileContent)
			if err != nil {
				return err
			}
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
