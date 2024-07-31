package internal

import (
	"rawdog-md/global"
	"rawdog-md/helper"
	"strings"

	"github.com/aymerick/raymond"
)

type Context map[string]any

func NewContexts(pages []Page) map[string]Context {
	pagesContext := make([]any, 0)
	for _, page := range pages {
		pageMap := make(map[string]any)

		// Remove leading slash from URL
		url := page.RelativeUrl
		if len(url) > 0 && url[0] == '/' {
			url = url[1:]
		}

		pageMap["$url"] = url
		pageMap["$type"] = strings.ToLower(page.Type.String())
		pageMap["$filename"] = helper.OmitFilenameExtension(page.Filename)
		pageMap["$body"] = page.Body

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
		url := page.RelativeUrl
		if len(url) > 0 && url[0] == '/' {
			url = url[1:]
		}

		context["$pages"] = pagesContext
		context["$buildMode"] = global.Config.BuildMode

		context["$url"] = url
		context["$type"] = strings.ToLower(page.Type.String())
		context["$filename"] = helper.OmitFilenameExtension(page.Filename)
		context["$body"] = raymond.SafeString(page.Body)

		if page.Frontmatter != nil {
			for k, v := range *page.Frontmatter {
				context[k] = v
			}
		}

		contexts[page.SourceAbsolutePath] = context
	}

	return contexts
}
