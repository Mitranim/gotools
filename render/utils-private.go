package render

import (
	// Standard
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	// Third party
	"github.com/Mitranim/gotools/utils"
)

/********************************* Constants *********************************/

// Errors used in this package. Each error starts with an http status code that
// can be converted to int with ErrorCode(err).
const (
	err404    = utils.Error("404 template not found")
	err500    = utils.Error("500 template rendering error")
	err500ISE = utils.Error("500 internal server error")
)

/********************************** Render ***********************************/

// Renders the given template at the given path or returns an error.
func renderAt(temp *template.Template, path string, data map[string]interface{}) ([]byte, error) {
	// Adjust and validate path.
	path, err := parsePath(temp, path)
	if err != nil {
		return nil, err
	}

	// Check for nil map.
	if data == nil {
		data = map[string]interface{}{}
	}

	// Mark path in data.
	if str, _ := data["path"].(string); str == "" {
		data["path"] = path
	}

	wr := new(utils.WR)
	err = temp.ExecuteTemplate(wr, path, data)
	if err != nil {
		return nil, err
	}

	return []byte(*wr), nil
}

/********************** Template Registration Utilities **********************/

// Traverses the given template directory and parses the files, creating
// templates under paths mimicking their path relative to that directory.
// Returns an error if anything goes wrong.
func readTemplates(dir string, temp *template.Template) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info == nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// Get filename without extension
		name := info.Name()
		name = strings.Replace(name, filepath.Ext(name), "", 1)

		// Get path without extension
		modpath := filepath.Join(filepath.Dir(path), name)

		// Remove prefix from path
		modpath = dePrefix(dir, modpath)

		// Get file contents
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return utils.Error(fmt.Sprintf("couldn't read file at path: %s, error: %#v\n", path, err))
		}

		// Parse template
		_, err = temp.New(modpath).Parse(string(bytes))
		if err != nil {
			return utils.Error(fmt.Sprintf("couldn't parse template at path: %s, error: %#v\n", modpath, err))
		}

		return nil
	})
}

// Traverses the inline file directory and reads each file into memory for
// future inlining. Returns an error if anything goes wrong.
func readInline(dir string, inlineFiles map[string]template.HTML) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		// Read file into memory.
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return utils.Error(fmt.Sprintf("couldn't read file at path: %s, error: %#v\n", path, err))
		}

		// Remove directory prefix from path.
		path = dePrefix(dir, path)

		// Put into map.
		inlineFiles[path] = template.HTML(bytes)

		return nil
	})
}

/****************************** Path Utilities *******************************/

// Takes a path to a resource and returns a reverse template hierarchy.
func pathsToTemplates(path string) []string {
	return reverse(makeTemplateHierarchy(path))
}

// Takes a path to a page and returns a slice of paths to consequtive
// hierarchical templates, where directory roots correspond to index
// templates. The original path comes last. Example:
//   blog/posts/my-post ->
//   []string{"index", "blog/index", "blog/posts/index", "blog/posts/my-post"}
// If the file name (the last name) is literally named `index`, it's skipped;
// its directory implicitly mandates this template.
func makeTemplateHierarchy(path string) []string {
	names := split(path)

	// Start at the implicit rootmost `index` layout and build the path list.
	paths := []string{"index"}

	if len(names) == 0 {
		return paths
	}

	// Loop over the directory names (in other words, all names except the last)
	// and build path names, as illustrated above.
	for index := range names[:len(names)-1] {
		dir := strings.Join(names[:index+1], "/")
		path := dir + "/index"
		paths = append(paths, path)
	}

	// If the file name is not `index`, add its path.
	if names[len(names)-1] != "index" {
		paths = append(paths, names[len(names)-1])
	}

	return paths
}

// Reverses a slice of strings.
func reverse(original []string) (reversed []string) {
	for i := len(original) - 1; i >= 0; i-- {
		reversed = append(reversed, original[i])
	}
	return reversed
}

// Removes the given prefix from a path. If the path isn't longer than the
// prefix, it's replaced with "/".
func dePrefix(prefix, path string) string {
	if len(path) <= len(prefix) {
		path = "/"
	} else {
		path = path[len(prefix)-1:]
	}
	return path
}

// Adjusts the path by dropping starting and ending slashes and checks if the
// path exists in the given template.
func parsePath(temp *template.Template, path string) (string, error) {
	// Drop starting slash, if any.
	if len(path) > 0 && path[0] == '/' {
		path = path[0:]
	}

	// Drop ending slash, if any.
	if len(path) > 0 && path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}

	// Make sure the path is still non-zero and the template exists.
	if len(path) == 0 || temp.Lookup(path) == nil {
		return path, err404
	}

	return path, nil
}

// Clears a slice of strings from empty strings.
func purify(paths []string) []string {
	result := []string{}
	for _, value := range paths {
		if value != "" {
			result = append(result, value)
		}
	}
	return result
}

// Splits the given path into a slice of strings, removing empty strings.
func split(path string) []string {
	return purify(strings.Split(path, "/"))
}

/*********************************** Other ***********************************/

// Returns an inlined file at the given path (if available) or an empty string,
// registering it in the given data map. Further calls with the same path and
// data map return an empty string.
func inline(state *StateInstance, path string, data map[string]interface{}) template.HTML {
	// Make sure we have an inline cache.
	cache, _ := data["inlined"].(map[string]bool)
	if cache == nil {
		cache = map[string]bool{}
	}

	// Check if we're in a development environment. If true, re-read the file from
	// the disk.
	if state.config.DevChecker != nil && state.config.DevChecker() {
		bytes, err := ioutil.ReadFile(state.config.InlineDir + "/" + path)
		if err == nil {
			state.inlineFiles[path] = template.HTML(bytes)
		}
	}

	// If it's already been inlined or if there's no such file, return an empty
	// string.
	str, ok := state.inlineFiles[path]
	if cache[path] || !ok {
		return ""
	}

	// Register and inline the file.
	cache[path] = true
	data["inlined"] = cache
	return str
}

// Logs stuff using a logger from a config, if any.
func log(state *StateInstance, values ...interface{}) {
	if state.config.Logger != nil {
		state.config.Logger(values...)
	}
}

// Converts the given error to a template path using a CodePath func passed in a
// config, if any. If it's omitted, uses a direct int to string conversion: 404
// -> "404".
func errorPath(state *StateInstance, err error) string {
	code := ErrorCode(err)
	if state.config.CodePath != nil {
		return state.config.CodePath(code)
	}
	return CodePath(code)
}
