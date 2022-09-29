package lang

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"goyave.dev/goyave/v4/config"
	"goyave.dev/goyave/v4/util/fsutil"
	"goyave.dev/goyave/v4/util/httputil"
)

var languages map[string]*Language
var mutex = &sync.RWMutex{}

// LoadDefault load the fallback language ("en-US").
// This function is intended for internal use only.
func LoadDefault() {
	mutex.Lock()
	defer mutex.Unlock()
	languages = make(map[string]*Language, 1)
	languages[enUS.name] = enUS.clone()
}

// LoadAllAvailableLanguages loads every language directory
// in the "resources/lang" directory if it exists.
func LoadAllAvailableLanguages() {
	mutex.Lock()
	defer mutex.Unlock()
	sep := string(os.PathSeparator)
	workingDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	langDirectory := workingDir + sep + "resources" + sep + "lang" + sep
	if fsutil.IsDirectory(langDirectory) {
		files, err := os.ReadDir(langDirectory)
		if err != nil {
			panic(err)
		}

		for _, f := range files {
			if f.IsDir() {
				load(f.Name(), langDirectory+sep+f.Name())
			}
		}
	}
}

// Load a language directory.
//
// Directory structure of a language directory:
//
//	en-UK
//	  ├─ locale.json (contains the normal language lines)
//	  ├─ rules.json  (contains the validation messages)
//	  └─ fields.json (contains the attribute-specific validation messages)
//
// Each file is optional.
func Load(language, path string) {
	mutex.Lock()
	defer mutex.Unlock()
	if fsutil.IsDirectory(path) {
		load(language, path)
	} else {
		panic(fmt.Sprintf("Failed loading language \"%s\", directory \"%s\" doesn't exist", language, path))
	}
}

func load(lang string, path string) {
	langStruct := &Language{name: lang}
	pathPrefix := path + string(os.PathSeparator)
	readLangFile(pathPrefix+"locale.json", &langStruct.lines)
	readLangFile(pathPrefix+"rules.json", &langStruct.validation.rules)
	readLangFile(pathPrefix+"fields.json", &langStruct.validation.fields)

	if existingLang, exists := languages[lang]; exists {
		mergeLang(existingLang, langStruct)
	} else {
		languages[lang] = langStruct
	}
}

// Get a language line.
//
// For validation rules messages and field names, use a dot-separated path:
// - "validation.rules.<rule_name>"
// - "validation.fields.<field_name>"
// For normal lines, just use the name of the line. Note that if you have
// a line called "validation", it won't conflict with the dot-separated paths.
//
// If not found, returns the exact "line" argument.
//
// The placeholders parameter is a variadic associative slice of placeholders and their
// replacement. In the following example, the placeholder ":username" will be replaced
// with the Name field in the user struct.
//
//	lang.Get("en-US", "greetings", ":username", user.Name)
func Get(lang string, line string, placeholders ...string) string {
	if !IsAvailable(lang) {
		return line
	}

	mutex.RLock()
	defer mutex.RUnlock()
	if strings.Count(line, ".") > 0 {
		path := strings.Split(line, ".")
		if path[0] == "validation" {
			switch path[1] {
			case "rules":
				if len(path) < 3 {
					return line
				}
				return convertEmptyLine(line, languages[lang].validation.rules[strings.Join(path[2:], ".")], placeholders)
			case "fields":
				if len(path) != 3 {
					return line
				}
				return convertEmptyLine(line, languages[lang].validation.fields[path[2]], placeholders)
			default:
				return line
			}
		}
	}

	return convertEmptyLine(line, languages[lang].lines[line], placeholders)
}

// IsAvailable returns true if the language is available.
func IsAvailable(lang string) bool {
	mutex.RLock()
	defer mutex.RUnlock()
	_, exists := languages[lang]
	return exists
}

// GetAvailableLanguages returns a slice of all loaded languages.
// This can be used to generate different routes for all languages
// supported by your applications.
//
//	/en/products
//	/fr/produits
//	...
func GetAvailableLanguages() []string {
	mutex.RLock()
	defer mutex.RUnlock()
	langs := []string{}
	for lang := range languages {
		langs = append(langs, lang)
	}
	return langs
}

// DetectLanguage detects the language to use based on the given lang string.
// The given lang string can use the HTTP "Accept-Language" header format.
//
// If "*" is provided, the default language will be used.
// If multiple languages are given, the first available language will be used,
// and if none are available, the default language will be used.
// If no variant is given (for example "en"), the first available variant will be used.
// For example, if "en-US" and "en-UK" are available and the request accepts "en",
// "en-US" will be used.
func DetectLanguage(lang string) string {
	values := httputil.ParseMultiValuesHeader(lang)
	for _, l := range values {
		if l.Value == "*" { // Accept anything, so return default language
			break
		}
		if IsAvailable(l.Value) {
			return l.Value
		}
		for key := range languages {
			if strings.HasPrefix(key, l.Value) {
				return key
			}
		}
	}

	return config.GetString("app.defaultLanguage")
}
