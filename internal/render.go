package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"

	"github.com/araddon/dateparse"
	"github.com/aymerick/raymond"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func init() {
	// create a custom helper for raymond for basic array methods, like reverse, sort, slice, etc.
	raymond.RegisterHelper("slice", func(context []interface{}, start int, end int, options *raymond.Options) string {
		result := options.Inverse()
		if !raymond.IsTrue(context) {
			return result
		}

		val := reflect.ValueOf(context)
		if val.Kind() != reflect.Array && val.Kind() != reflect.Slice {
			fmt.Println("WARNING: Context must be an array or slice")
			return result
		}

		if end < 0 {
			fmt.Println("WARNING: End index must be greater than or equal to 0")
			return result
		}

		if start > end {
			fmt.Println("WARNING: Start index must be less than or equal to end index")
			return result
		}

		if start == end {
			return result
		}

		if start >= val.Len() {
			fmt.Println("WARNING: Start index must be less than the length of the array or slice")
			return result
		}

		start = max(start, 0)
		end = min(end, val.Len())

		ctxSlice := context[start:end]

		return options.FnWith(ctxSlice)
	})

	raymond.RegisterHelper("withSortByDate", func(context []interface{}, key string, order string, options *raymond.Options) string {
		result := options.Inverse()
		if !raymond.IsTrue(context) {
			return result
		}

		val := reflect.ValueOf(context)
		if val.Kind() != reflect.Array && val.Kind() != reflect.Slice {
			fmt.Println("WARNING: Context must be an array or slice")
			return result
		}

		if val.Len() == 0 {
			return result
		}

		if order != "asc" && order != "desc" {
			fmt.Println("WARNING: Order must be either 'asc' or 'desc'")
			return result
		}

		if key == "" {
			fmt.Println("WARNING: Key must be a non-empty string")
			return result
		}

		sort.Slice(context, func(i, j int) bool {
			iDateStr, ok := context[i].(map[string]interface{})[key].(string)
			if !ok {
				return false
			}

			jDateStr, ok := context[j].(map[string]interface{})[key].(string)
			if !ok {
				return false
			}

			iDate, err := dateparse.ParseAny(iDateStr)
			if err != nil {
				return false
			}

			jDate, err := dateparse.ParseAny(jDateStr)
			if err != nil {
				return false
			}

			if order == "desc" {
				return iDate.After(jDate)
			}

			return iDate.Before(jDate)
		})

		return options.FnWith(context)
	})

	raymond.RegisterHelper("withFilter", func(context []interface{}, key string, value interface{}, options *raymond.Options) string {
		result := options.Inverse()
		if !raymond.IsTrue(context) {
			return result
		}

		val := reflect.ValueOf(context)
		if val.Kind() != reflect.Array && val.Kind() != reflect.Slice {
			fmt.Println("WARNING: Context must be an array or slice")
			return result
		}

		if val.Len() == 0 {
			return result
		}

		if key == "" {
			fmt.Println("WARNING: Key must be a non-empty string")
			return result
		}

		filtered := make([]interface{}, 0)
		for _, item := range context {
			itemMap, ok := item.(map[string]interface{})
			if !ok {
				continue
			}

			if itemMap[key] == value {
				filtered = append(filtered, item)
			}
		}

		return options.FnWith(filtered)
	})

	raymond.RegisterHelper("toJson", func(context interface{}, options *raymond.Options) string {
		if context == nil {
			return ""
		}

		b, err := json.Marshal(context)
		if err != nil {
			fmt.Println("WARNING: Failed to marshal context to JSON")
			return ""
		}

		return string(b)
	})

	raymond.RegisterHelper("toJsonDiv", func(context interface{}, options *raymond.Options) raymond.SafeString {
		if context == nil {
			return ""
		}

		b, err := json.MarshalIndent(context, "", "\t")
		if err != nil {
			fmt.Println("WARNING: Failed to marshal context to pretty JSON")
			return ""
		}

		result := fmt.Sprintf("<div style=\"white-space: pre-wrap;\">%s</div>", string(b))

		return raymond.SafeString(result)
	})

	raymond.RegisterHelper("logJson", func(context interface{}, options *raymond.Options) string {
		if context == nil {
			return ""
		}

		b, err := json.MarshalIndent(context, "", "\t")
		if err != nil {
			fmt.Println("WARNING: Failed to marshal context to pretty JSON")
			return ""
		}

		fmt.Println(string(b))

		return ""
	})
}

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
