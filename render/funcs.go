package render

import (
	"html/template"
)

// Generates a map of base template funcs that closure a reference to a state
// object.
func makeTemplateFuncs(state *StateInstance) template.FuncMap {
	return template.FuncMap{

		// Modifies the title of the data object, appending the given string. Always
		// returns an empty string.
		"title": func(title string, data map[string]interface{}) (result string) {
			str, _ := data["title"].(string)
			if str == "" {
				data["title"] = title
			} else {
				str += " | " + title
				data["title"] = str
			}
			return
		},

		// Includes the given standalone template only once during the lifetime of the
		// &data.
		"import": func(path string, data map[string]interface{}) template.HTML {
			// Make sure we have an import cache.
			cache, _ := data["imported"].(map[string]bool)
			if cache == nil {
				cache = map[string]bool{}
			}
			// If it's already been imported, return an empty string.
			if cache[path] {
				return ""
			}
			// Otherwise register and import it.
			cache[path] = true
			data["imported"] = cache
			bytes, err := state.RenderStandalone(path, data)
			if err != nil {
				return ""
			}
			return template.HTML(bytes)
		},

		// Inlines the given file only once during the lifetime of the &data.
		"inline": func(path string, data map[string]interface{}) template.HTML {
			return inline(state, path, data)
		},

		// Inlines the given file as a stylesheet.
		"inlineStyle": func(path string, data map[string]interface{}) template.HTML {
			text := inline(state, path, data)
			if text == "" {
				return ""
			}
			return "<style>\n" + text + "\n</style>"
		},

		// Inlines the given file as a script.
		"inlineScript": func(path string, data map[string]interface{}) template.HTML {
			text := inline(state, path, data)
			if text == "" {
				return ""
			}
			return `<script type="text/javascript">` + "\n" + text + "\n" + `</script>`
		},
	}
}
