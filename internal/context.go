package internal

import (
	"strings"

	"rawd/global"
	"rawd/helper"

	"github.com/aymerick/raymond"
)

type Context map[string]any

func NewContexts(pages []Page) map[string]Context {
	pagesContext := make([]any, 0)
	for _, page := range pages {
		pageMap := make(map[string]any)

		// Remove leading slash from URL
		path := page.RelativeUrl
		if len(path) > 0 && path[0] == '/' {
			path = path[1:]
		}

		pageMap["Path"] = path
		pageMap["Type"] = strings.ToLower(page.Type.String())
		pageMap["Filename"] = helper.OmitFilenameExtension(page.Filename)
		pageMap["Body"] = page.Body

		if page.Frontmatter != nil {
			for k, v := range *page.Frontmatter {
				pageMap[k] = v
			}
		}

		pagesContext = append(pagesContext, pageMap)
	}

	contexts := make(map[string]Context)
	for _, page := range pages {
		context := make(Context)

		// Remove leading slash from URL
		path := page.RelativeUrl
		if len(path) > 0 && path[0] == '/' {
			path = path[1:]
		}

		context["Pages"] = pagesContext
		context["BuildMode"] = global.Config.BuildMode

		context["Path"] = path
		context["Type"] = strings.ToLower(page.Type.String())
		context["Filename"] = helper.OmitFilenameExtension(page.Filename)
		context["Body"] = raymond.SafeString(page.Body)

		if page.Frontmatter != nil {
			for k, v := range *page.Frontmatter {
				context[k] = v
			}
		}

		contexts[page.SourceAbsolutePath] = context
	}

	return contexts
}
