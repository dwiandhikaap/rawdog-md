package internal

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/araddon/dateparse"
	"github.com/dwiandhikaap/rawdog-md/global"
	"github.com/yuin/goldmark"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	enclave "github.com/quail-ink/goldmark-enclave"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	anchor "go.abhg.dev/goldmark/anchor"

	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
)

func createMarkdownParser() goldmark.Markdown {
	extensions := []goldmark.Extender{}

	// GFM
	if global.Config.UserConfig.MarkdownPlugins.GFM.Enabled {
		extensions = append(extensions, extension.GFM)
	}

	// CJK
	if global.Config.UserConfig.MarkdownPlugins.CJK.Enabled {
		extensions = append(extensions, extension.CJK)
	}

	// DefinitionList
	if global.Config.UserConfig.MarkdownPlugins.DefinitionList.Enabled {
		extensions = append(extensions, extension.DefinitionList)
	}

	// Footnote
	if global.Config.UserConfig.MarkdownPlugins.Footnote.Enabled {
		extensions = append(extensions, extension.Footnote)
	}

	// Table
	if global.Config.UserConfig.MarkdownPlugins.Table.Enabled {
		extensions = append(extensions, extension.Table)
	}

	// Strikethrough
	if global.Config.UserConfig.MarkdownPlugins.Strikethrough.Enabled {
		extensions = append(extensions, extension.Strikethrough)
	}

	// Typographer
	if global.Config.UserConfig.MarkdownPlugins.Typographer.Enabled {
		extensions = append(extensions, extension.Typographer)
	}

	// TaskList
	if global.Config.UserConfig.MarkdownPlugins.TaskList.Enabled {
		extensions = append(extensions, extension.TaskList)
	}

	// Linkify
	if global.Config.UserConfig.MarkdownPlugins.Linkify.Enabled {
		extensions = append(extensions, extension.Linkify)
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

	rendererOptions := []renderer.Option{
		html.WithHardWraps(),
	}

	if global.Config.UserConfig.Options.Html.Unsafe {
		rendererOptions = append(rendererOptions, html.WithUnsafe())
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
			rendererOptions...,
		),
	)
}

var md goldmark.Markdown

func init() {
	InitMarkdownParser()
}

func InitMarkdownParser() {
	md = createMarkdownParser()
}

func convertMarkdown(content string) (string, error) {
	var buf bytes.Buffer
	if err := md.Convert([]byte(content), &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func getField(obj any, fieldName string) (reflect.Value, bool) {
	var val reflect.Value

	defer func() {
		if r := recover(); r != nil {
			val = reflect.Value{}
		}
	}()

	val = reflect.ValueOf(obj)

	if val.Kind() == reflect.Map {
		for _, key := range val.MapKeys() {
			if key.String() == fieldName {
				mapVal := val.MapIndex(key)
				if !mapVal.IsValid() {
					return reflect.Value{}, false
				}
				return mapVal, true
			}
		}
		return reflect.Value{}, false
	}

	return reflect.Value{}, false
}

func sortDate(values []any, order string) []any {
	order = strings.ToLower(order)
	if order != "asc" && order != "desc" {
		fmt.Println("sorting order should be either 'asc' or 'desc'")
	}

	compareInner := func(i, j int) bool {
		fieldI, hasFieldI := getField(values[i], "Date")
		fieldJ, hasFieldJ := getField(values[j], "Date")

		if !hasFieldI || !hasFieldJ {
			return false
		}

		if fieldI.Kind() == reflect.Interface {
			stringTimeI := fieldI.Interface().(string)
			stringTimeJ := fieldJ.Interface().(string)

			timeI, errI := dateparse.ParseAny(stringTimeI)
			timeJ, errJ := dateparse.ParseAny(stringTimeJ)

			if errI != nil || errJ != nil {
				return false
			}

			return timeI.Before(timeJ)
		}

		return false
	}

	compare := func(i, j int) bool {
		if order == "desc" {
			return !compareInner(i, j)
		}
		return compareInner(i, j)
	}

	sort.SliceStable(values, compare)

	return values
}

func renderTextTemplate(content string, context map[string]interface{}) (string, error) {
	funcMap := sprig.FuncMap()

	funcMap["sortDate"] = sortDate

	tmpl, err := template.New("text").Funcs(funcMap).Parse(content)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, context); err != nil {
		return "", err
	}

	return buf.String(), nil
}
