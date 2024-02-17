package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"rawdog-md/global"
	"rawdog-md/helper"
	"strings"

	"github.com/adrg/frontmatter"
)

type Frontmatter map[string]any
type PageType int

const (
	Markdown PageType = iota
	Handlebars
	Html
)

func (p PageType) String() string {
	return [...]string{"Markdown", "Handlebars", "Html"}[p]
}

type Page struct {
	SourceRelativePath string
	SourceAbsolutePath string
	Filename           string
	RelativeUrl        string
	Type               PageType

	Frontmatter *Frontmatter

	TemplateName *string
	Template     *Template

	Body   string
	Output string
}

func (p *Page) Render(contextMap map[string]Context) error {
	context := contextMap[p.SourceAbsolutePath]

	if p.Type == Html {
		p.Output = p.Body
		return nil
	}

	if p.Type == Markdown {
		html, err := renderHandlebars(p.Template.Content, context)
		if err != nil {
			return err
		}

		p.Output = html
		return nil
	}

	if p.Type == Handlebars {
		html, err := renderHandlebars(p.Template.Content, context)
		if err != nil {
			return err
		}

		p.Output = html
		return nil
	}

	return nil
}

func (p *Page) Dump() {
	template := "False"
	if p.Template != nil {
		template = p.Template.AbsolutePath
	}

	dump := fmt.Sprintf(`
Page: '%s'
Type: %s
HasTemplate: %v
RelativeUrl: %s
Body: %s
Output: %s
	`,
		p.SourceAbsolutePath,
		p.Type.String(),
		template,
		p.RelativeUrl,
		helper.TruncateString(p.Body, 30),
		helper.TruncateString(p.Output, 400),
	)

	fmt.Println(dump)
}

func NewPage(absolutePath string) (*Page, error) {
	relativePath, err := filepath.Rel(global.Config.RootAbsolutePath, absolutePath)
	filename := filepath.Base(absolutePath)

	if err != nil {
		return nil, err
	}

	Page := &Page{
		Type:               Markdown,
		Filename:           filepath.Base(relativePath),
		SourceRelativePath: relativePath,
		SourceAbsolutePath: absolutePath,
	}

	if filepath.Ext(filename) == ".md" {
		err = loadMarkdown(Page)
	} else if filepath.Ext(filename) == ".html" {
		err = loadHtml(Page)
	} else if filepath.Ext(filename) == ".hbs" {
		err = loadHandlebars(Page)
	} else {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return Page, nil
}

func loadMarkdown(p *Page) error {
	fileContent, err := os.ReadFile(p.SourceAbsolutePath)
	if err != nil {
		return err
	}

	var fm Frontmatter

	reader := strings.NewReader(string(fileContent))
	content, err := frontmatter.Parse(reader, &fm)
	if err != nil {
		return err
	}

	err = validateFrontmatter(fm)
	if err != nil {
		return err
	}

	p.Frontmatter = &fm

	templateName := fm["template"].(string) + ".hbs"
	p.TemplateName = &templateName

	relativePath := p.SourceRelativePath
	relativePageWithoutExt := relativePath[:len(relativePath)-len(filepath.Ext(p.Filename))]
	p.RelativeUrl = relativePageWithoutExt + ".html"
	p.RelativeUrl = strings.Replace(p.RelativeUrl, "\\", "/", -1)
	p.RelativeUrl = p.RelativeUrl[5:]

	markdownHtml, err := convertMarkdown(string(content))
	if err != nil {
		return err
	}
	p.Body = markdownHtml
	return nil
}

func loadHtml(p *Page) error {
	fileContent, err := os.ReadFile(p.SourceAbsolutePath)
	if err != nil {
		return err
	}

	p.Body = string(fileContent)
	p.RelativeUrl = p.SourceRelativePath
	p.RelativeUrl = strings.Replace(p.RelativeUrl, "\\", "/", -1)
	p.RelativeUrl = p.RelativeUrl[5:]

	return nil
}

func loadHandlebars(p *Page) error {
	fileContent, err := os.ReadFile(p.SourceAbsolutePath)
	if err != nil {
		return err
	}

	p.Template = &Template{
		AbsolutePath: p.SourceAbsolutePath,
		Filename:     p.Filename,
		Content:      string(fileContent),
	}
	p.Body = ""

	relativePageWithoutExt := p.SourceRelativePath[:len(p.SourceRelativePath)-len(filepath.Ext(p.Filename))]
	p.RelativeUrl = relativePageWithoutExt + ".html"
	p.RelativeUrl = strings.Replace(p.RelativeUrl, "\\", "/", -1)
	p.RelativeUrl = p.RelativeUrl[5:]

	return nil
}

func validateFrontmatter(frontmatter Frontmatter) error {
	if len(frontmatter) == 0 {
		return fmt.Errorf("markdown file must contain frontmatter")
	}

	if _, ok := frontmatter["template"]; !ok {
		return fmt.Errorf("frontmatter must contain a 'template' field")
	}

	return nil
}
