package internal

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/dwiandhikaap/rawdog-md/global"
	"github.com/yuin/goldmark"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	enclave "github.com/quail-ink/goldmark-enclave"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	anchor "go.abhg.dev/goldmark/anchor"

	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func createMarkdownParser() goldmark.Markdown {
	extensions := []goldmark.Extender{
		extension.GFM,
		extension.CJK,
		extension.DefinitionList,
		extension.Footnote,
		extension.Table,
		extension.Strikethrough,
		extension.Typographer,
		extension.TaskList,
		extension.Linkify,
	}

	// Highlighting
	if global.Config.UserConfig.MarkdownPlugins.Highlighting.Enabled {
		formatOpts := []chromahtml.Option{}
		highlightingOpts := []highlighting.Option{highlighting.WithGuessLanguage(true)}

		if global.Config.UserConfig.MarkdownPlugins.Highlighting.Style == nil {
			formatOpts = append(formatOpts, chromahtml.WithClasses(true))
		} else {
			highlightingOpts = append(highlightingOpts, highlighting.WithStyle(*global.Config.UserConfig.MarkdownPlugins.Highlighting.Style))
		}

		if global.Config.UserConfig.MarkdownPlugins.Highlighting.UseLineNumbers {
			formatOpts = append(formatOpts, chromahtml.WithLineNumbers(true))
		}

		highlightingOpts = append(highlightingOpts, highlighting.WithFormatOptions(formatOpts...))

		extensions = append(extensions, highlighting.NewHighlighting(highlightingOpts...))
	}

	// Enclave
	if global.Config.UserConfig.MarkdownPlugins.Enclave.Enabled {
		extensions = append(extensions, enclave.New(&enclave.Config{}))
	}

	// Anchor
	if global.Config.UserConfig.MarkdownPlugins.Anchor.Enabled {
		position := anchor.Before
		if global.Config.UserConfig.MarkdownPlugins.Anchor.Position == "right" {
			position = anchor.After
		}

		extensions = append(extensions, &anchor.Extender{
			Position: position,
			Texter:   anchor.Text(global.Config.UserConfig.MarkdownPlugins.Anchor.Text),
			Attributer: anchor.Attributes{
				"class": global.Config.UserConfig.MarkdownPlugins.Anchor.Class,
			},
		})
	}

	return goldmark.New(
		goldmark.WithExtensions(
			extensions...,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithHeadingAttribute(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
}

func convertMarkdown(content string) (string, error) {
	md := createMarkdownParser()

	var buf bytes.Buffer
	if err := md.Convert([]byte(content), &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func renderTextTemplate(content string, context map[string]interface{}) (string, error) {
	tmpl, err := template.New("text").Funcs(sprig.FuncMap()).Parse(content)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, context); err != nil {
		return "", err
	}

	return buf.String(), nil
}
