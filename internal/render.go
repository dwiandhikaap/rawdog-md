package internal

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/yuin/goldmark"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	enclave "github.com/quail-ink/goldmark-enclave"
	highlighting "github.com/yuin/goldmark-highlighting/v2"

	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

var md = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		extension.CJK,
		extension.DefinitionList,
		extension.Footnote,
		extension.Table,
		extension.Strikethrough,
		extension.Typographer,
		extension.TaskList,
		extension.Linkify,
		highlighting.NewHighlighting(
			highlighting.WithFormatOptions(
				chromahtml.WithLineNumbers(true),
				chromahtml.WithClasses(true),
			),
			highlighting.WithGuessLanguage(true),
		),
		enclave.New(&enclave.Config{
			DefaultImageAltPrefix: "imggg",
		}),
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
	goldmark.WithRendererOptions(
		html.WithHardWraps(),
		html.WithXHTML(),
	),
)

func convertMarkdown(content string) (string, error) {
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
