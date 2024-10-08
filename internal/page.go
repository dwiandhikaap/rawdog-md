package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dwiandhikaap/rawdog-md/global"
	"github.com/dwiandhikaap/rawdog-md/helper"

	"github.com/adrg/frontmatter"
	"gopkg.in/yaml.v2"
)

type Frontmatter map[string]any
type PageType int

const (
	Markdown PageType = iota
	TextTemplate
	Html
)

func (p PageType) String() string {
	return [...]string{"Markdown", "Html", "TextTemplate"}[p]
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

	Body   string // Rendered body, ready to be used in a template as Body
	Output string // Final output, ready to write to file
}

func NewPage(absolutePath string) (*Page, error) {
	relativePath, err := filepath.Rel(global.Config.RootAbsolutePath, absolutePath)
	filename := filepath.Base(absolutePath)

	if err != nil {
		return nil, err
	}

	Page := &Page{
		Filename:           filepath.Base(relativePath),
		SourceRelativePath: relativePath,
		SourceAbsolutePath: absolutePath,
	}

	if filepath.Ext(filename) == ".md" {
		err = loadMarkdown(Page)
		Page.Type = Markdown
	} else if filepath.Ext(filename) == ".html" {
		err = loadHtml(Page)
		Page.Type = Html
	} else if filepath.Ext(filename) == ".templ" {
		err = loadTextTemplate(Page)
		Page.Type = TextTemplate
	} else {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return Page, nil
}

func (p *Page) Reload() error {
	var err error
	if filepath.Ext(p.Filename) == ".md" {
		err = loadMarkdown(p)
	} else if filepath.Ext(p.Filename) == ".html" {
		err = loadHtml(p)
	} else if filepath.Ext(p.Filename) == ".templ" {
		err = loadTextTemplate(p)
	}

	if err != nil {
		return fmt.Errorf("error reloading page '%s': %v", p.SourceAbsolutePath, err)
	}

	return nil
}

func (p *Page) Render(contextMap map[string]Context) error {
	context := contextMap[p.SourceAbsolutePath]

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("PANIC: unable to render '%s':\n%v\n", p.SourceAbsolutePath, r)
		}
	}()

	if p.Type == Html {
		p.Output = p.Body
		if global.Config.UserConfig.Options.Minify.HTML {
			out, err := helper.Minifier.Bytes("text/html", []byte(p.Output))
			if err != nil {
				return err
			}

			p.Output = string(out)

			return nil
		}
		return nil
	}

	if p.Type == Markdown {
		if p.Template == nil {
			return fmt.Errorf("page '%s' has no template", p.SourceAbsolutePath)
		}

		html, err := renderTextTemplate(p.Template.Content, context)
		if err != nil {
			return err
		}

		if global.Config.UserConfig.Options.Minify.HTML {
			out, err := helper.Minifier.Bytes("text/html", []byte(html))
			if err != nil {
				return err
			}

			html = string(out)
		}

		p.Output = html
		return nil
	}

	if p.Type == TextTemplate {
		html, err := renderTextTemplate(p.Template.Content, context)
		if err != nil {
			return err
		}

		if global.Config.UserConfig.Options.Minify.HTML {
			out, err := helper.Minifier.Bytes("text/html", []byte(html))
			if err != nil {
				return err
			}

			p.Output = string(out)
		}

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

	templateName := fm["Template"].(string) + ".templ"
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

	format := frontmatter.NewFormat("<!-- fm-yaml-start", "fm-yaml-end -->", yaml.Unmarshal)

	var fm Frontmatter
	reader := strings.NewReader(string(fileContent))
	_, err = frontmatter.Parse(reader, &fm, format)
	if err != nil {
		return err
	}

	p.Body = string(fileContent)
	p.RelativeUrl = p.SourceRelativePath
	p.RelativeUrl = strings.Replace(p.RelativeUrl, "\\", "/", -1)
	p.RelativeUrl = p.RelativeUrl[5:]
	p.Frontmatter = &fm

	return nil
}

func loadTextTemplate(p *Page) error {
	fileContent, err := os.ReadFile(p.SourceAbsolutePath)
	if err != nil {
		return err
	}

	p.Template = &Template{
		AbsolutePath: p.SourceAbsolutePath,
		Filename:     p.Filename,
		Content:      string(fileContent),
	}
	p.Body = "" // Text templates have no body as it will be treated as a template under the hood.. i guess..

	var fm Frontmatter

	format := frontmatter.NewFormat("{{/* fm-yaml-start", "fm-yaml-end */}}", yaml.Unmarshal)

	reader := strings.NewReader(string(fileContent))
	_, err = frontmatter.Parse(reader, &fm, format)
	if err != nil {
		return err
	}

	relativePageWithoutExt := p.SourceRelativePath[:len(p.SourceRelativePath)-len(filepath.Ext(p.Filename))]
	p.RelativeUrl = relativePageWithoutExt + ".html"
	p.RelativeUrl = strings.Replace(p.RelativeUrl, "\\", "/", -1)
	p.RelativeUrl = p.RelativeUrl[5:]
	p.Frontmatter = &fm

	return nil
}

func validateFrontmatter(frontmatter Frontmatter) error {
	if len(frontmatter) == 0 {
		return fmt.Errorf("markdown file must contain frontmatter")
	}

	if _, ok := frontmatter["Template"]; !ok {
		return fmt.Errorf("frontmatter must contain a 'Template' field")
	}

	return nil
}
