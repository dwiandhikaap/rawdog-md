package internal

import (
	"rawdog-md/helper"
)

type Context map[string]any

func NewContexts(pages []Page) map[string]Context {
	pagesContext := make([]any, 0)
	for _, page := range pages {
		pageMap := make(map[string]any)

		pageMap["_url"] = page.RelativeUrl
		pageMap["_type"] = page.Type.String()
		pageMap["_filename"] = helper.OmitFilenameExtension(page.Filename)
		pageMap["_body"] = page.Body

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
		context["_pages"] = pagesContext
		context["_url"] = page.RelativeUrl
		context["_type"] = page.Type.String()
		context["_filename"] = helper.OmitFilenameExtension(page.Filename)
		context["_body"] = page.Body

		if page.Frontmatter != nil {
			for k, v := range *page.Frontmatter {
				context[k] = v
			}
		}

		contexts[page.SourceAbsolutePath] = context
	}

	return contexts
}
