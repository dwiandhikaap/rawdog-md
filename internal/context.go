package internal

import (
	"rawdog-md/helper"
)

type Context map[string]any

func NewContexts(pages []Page) map[string]Context {
	pagesContext := make([]any, 0)
	for _, page := range pages {
		pageMap := make(map[string]any)

		pageMap["$url"] = page.RelativeUrl
		pageMap["$type"] = page.Type.String()
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
		context["$pages"] = pagesContext
		context["$url"] = page.RelativeUrl
		context["$type"] = page.Type.String()
		context["$filename"] = helper.OmitFilenameExtension(page.Filename)
		context["$body"] = page.Body

		if page.Frontmatter != nil {
			for k, v := range *page.Frontmatter {
				context[k] = v
			}
		}

		contexts[page.SourceAbsolutePath] = context
	}

	return contexts
}
