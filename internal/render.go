package internal

import (
	"bytes"

	"github.com/aymerick/raymond"
	"github.com/yuin/goldmark"
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

func renderHandlebars(content string, context map[string]interface{}) (string, error) {
	out, err := raymond.Render(content, context)
	if err != nil {
		return "", err
	}

	return out, nil
}
