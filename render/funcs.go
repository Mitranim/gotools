package render

import (
	"html/template"
)

// Base template funcs.
var templateFuncs = template.FuncMap{

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
		bytes, err := RenderStandalone(path, data)
		if err != nil {
			return ""
		}
		return template.HTML(bytes)
	},

	// Inlines the given file only once during the lifetime of the &data.
	"inline": func(path string, data map[string]interface{}) template.HTML {
		return inline(path, data)
	},

	// Inlines the given file as a stylesheet.
	"inlineStyle": func(path string, data map[string]interface{}) template.HTML {
		style := "<style>\n" + inline(path, data) + "\n</style>"
		return style
	},

	// Inlines the given file as a script.
	"inlineScript": func(path string, data map[string]interface{}) template.HTML {
		script := `<script type="text/javascript">` + "\n" + inline(path, data) + "\n" + `</script>`
		return script
	},
}
