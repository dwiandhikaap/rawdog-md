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

	var Page *Page
	if filepath.Ext(filename) == ".md" {
		Page, err = newMarkdownPage(relativePath)
	} else if filepath.Ext(filename) == ".html" {
		Page, err = newHtmlPage(relativePath)
	} else if filepath.Ext(filename) == ".hbs" {
		Page, err = newHandlebarsPage(relativePath)
	} else {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return Page, nil
}

func newMarkdownPage(relativePath string) (*Page, error) {
	absolutePath := filepath.Join(global.Config.RootAbsolutePath, relativePath)

	Page := &Page{
		Type:               Markdown,
		Filename:           filepath.Base(relativePath),
		SourceAbsolutePath: absolutePath,
	}

	fileContent, err := os.ReadFile(absolutePath)
	if err != nil {
		return nil, err
	}

	var fm Frontmatter

	reader := strings.NewReader(string(fileContent))
	content, err := frontmatter.Parse(reader, &fm)
	if err != nil {
		return nil, err
	}

	err = validateFrontmatter(fm)
	if err != nil {
		return nil, err
	}

	markdownHtml, err := convertMarkdown(string(content))
	if err != nil {
		return nil, err
	}

	Page.Frontmatter = &fm
	Page.Body = markdownHtml

	templateName := fm["template"].(string) + ".hbs"
	Page.TemplateName = &templateName

	relativePageWithoutExt := relativePath[:len(relativePath)-len(filepath.Ext(Page.Filename))]
	Page.RelativeUrl = relativePageWithoutExt + ".html"
	Page.RelativeUrl = strings.Replace(Page.RelativeUrl, "\\", "/", -1)
	Page.RelativeUrl = Page.RelativeUrl[5:]

	return Page, nil
}

func newHtmlPage(relativePath string) (*Page, error) {
	absolutePath := filepath.Join(global.Config.RootAbsolutePath, relativePath)

	Page := &Page{
		Type:               Html,
		Filename:           filepath.Base(relativePath),
		SourceAbsolutePath: absolutePath,
	}

	fileContent, err := os.ReadFile(absolutePath)
	if err != nil {
		return nil, err
	}

	Page.Body = string(fileContent)
	Page.RelativeUrl = relativePath
	Page.RelativeUrl = strings.Replace(Page.RelativeUrl, "\\", "/", -1)
	Page.RelativeUrl = Page.RelativeUrl[5:]

	return Page, nil
}

func newHandlebarsPage(relativePath string) (*Page, error) {
	absolutePath := filepath.Join(global.Config.RootAbsolutePath, relativePath)

	Page := &Page{
		Type:               Handlebars,
		Filename:           filepath.Base(relativePath),
		SourceAbsolutePath: absolutePath,
	}

	fileContent, err := os.ReadFile(absolutePath)
	if err != nil {
		return nil, err
	}

	Page.Template = &Template{
		AbsolutePath: absolutePath,
		Filename:     Page.Filename,
		Content:      string(fileContent),
	}
	Page.Body = ""

	relativePageWithoutExt := relativePath[:len(relativePath)-len(filepath.Ext(Page.Filename))]
	Page.RelativeUrl = relativePageWithoutExt + ".html"
	Page.RelativeUrl = strings.Replace(Page.RelativeUrl, "\\", "/", -1)
	Page.RelativeUrl = Page.RelativeUrl[5:]

	return Page, nil
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
