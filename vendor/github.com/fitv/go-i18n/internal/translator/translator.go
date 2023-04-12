package translator

import (
	"fmt"
	"strings"
)

type Translator struct {
	lang map[string]interface{}
}

// New returns a Translator instance
func New(lang map[string]interface{}) *Translator {
	return &Translator{lang: lang}
}

// Trans returns language translation by the given key
func (t *Translator) Trans(key string, args ...interface{}) string {
	value, ok := t.get(key)
	if !ok {
		return key
	}
	if len(args) == 0 {
		return value
	}

	dict, ok := args[0].(map[string]interface{})
	if !ok {
		return fmt.Sprintf(value, args...)
	}
	for key, val := range dict {
		value = strings.Replace(value, "{"+key+"}", fmt.Sprintf("%v", val), -1)
	}
	return value
}

// get returns language translation from the Translator
func (t *Translator) get(key string) (str string, exists bool) {
	source := t.lang
	keys := strings.Split(key, ".")
	last := len(keys) - 1

	for i, k := range keys {
		val, ok := source[k]
		if !ok {
			return
		}

		switch v := val.(type) {
		case string:
			if i == last {
				return v, true
			}
			return
		case map[string]interface{}:
			source = v
		default:
			return
		}
	}
	return
}
