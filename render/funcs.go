package render

import (
	// Standard
	"html/template"
)

// Generates a map of base template funcs that closure a reference to a state
// object.
func makeTemplateFuncs(state *stateInstance) template.FuncMap {
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

		// Includes the given template only once during the lifetime of the &data.
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
			bytes, err := state.RenderOne(path, data)
			if err != nil {
				return ""
			}
			return template.HTML(bytes)
		},

		// Checks if the given link is active. If true, returns "active", else "".
		"active": func(link string, data map[string]interface{}) string {
			return active(link, data)
		},

		// Same as "active" but also includes the class attribute.
		"act": func(link string, data map[string]interface{}) template.HTMLAttr {
			act := active(link, data)
			if len(act) == 0 {
				return ""
			}
			return "class=\"active\""
		},

		// Prints a background-image style with the given src.
		"bgImg": func(src string) template.HTMLAttr {
			return template.HTMLAttr(`style="background-image: url(/img/` + src + `)"`)
		},

		// Same as bgImg but without the /img prefix.
		"bgUrl": func(src string) template.HTMLAttr {
			return template.HTMLAttr(`style="background-image: url(` + src + `)"`)
		},

		// Makes an opening tag from the given string. May include additional
		// attributes.
		"tag": func(name string, attrs ...interface{}) template.HTML {
			str := join(attrs...)
			if str != "" {
				str = " " + str
			}
			return template.HTML("<" + name + str + ">")
		},

		// Makes a closing tag from the given string.
		"untag": func(name string) template.HTML {
			return template.HTML("</" + name + ">")
		},

		// Inlines the given file only once during the lifetime of the &data.
		"inline": func(path string, data map[string]interface{}) template.HTML {
			return inline(state, path, data)
		},

		// Inlines the given file as a stylesheet.
		"inlineStyle": func(path string, data map[string]interface{}) template.HTML {
			return inlineStyle(state, path, data)
		},

		// Inlines the given file as a script.
		"inlineScript": func(path string, data map[string]interface{}) template.HTML {
			return inlineScript(state, path, data)
		},

		// Inlines the given file only in production.
		"inlineProd": func(path string, data map[string]interface{}) template.HTML {
			if isDev(state) {
				return ""
			}
			return inline(state, path, data)
		},

		// Inlines the given file as a stylesheet only in production.
		"inlineStyleProd": func(path string, data map[string]interface{}) template.HTML {
			if isDev(state) {
				return ""
			}
			return inlineStyle(state, path, data)
		},

		// Inlines the given file as a script only in production.
		"inlineScriptProd": func(path string, data map[string]interface{}) template.HTML {
			if isDev(state) {
				return ""
			}
			return inlineScript(state, path, data)
		},
	}
}
