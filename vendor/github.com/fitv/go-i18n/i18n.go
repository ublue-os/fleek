package i18n

import (
	"embed"
	"strings"

	"github.com/fitv/go-i18n/internal/translator"
	"gopkg.in/yaml.v3"
)

var emptyTrans = translator.New(make(map[string]interface{}))

type I18n struct {
	defaultLocale string
	transMap      map[string]*translator.Translator
}

// New returns an I18n instance
func New(fs embed.FS, path string) (*I18n, error) {
	i18n := &I18n{
		transMap: make(map[string]*translator.Translator),
	}

	dirEntries, err := fs.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range dirEntries {
		lang := make(map[string]interface{})

		file, err := fs.ReadFile(path + "/" + entry.Name())
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(file, &lang)
		if err != nil {
			return nil, err
		}

		local := strings.Split(entry.Name(), ".")[0]
		i18n.transMap[local] = translator.New(lang)
	}
	return i18n, nil
}

// SetDefaultLocale set the default locale
func (i *I18n) SetDefaultLocale(local string) {
	i.defaultLocale = local
}

// Locale returns the translator instance by the given locale
func (i *I18n) Locale(locale string) *translator.Translator {
	trans, ok := i.transMap[locale]
	if ok {
		return trans
	}
	return emptyTrans
}

// Trans returns language translation by the given key
func (i *I18n) Trans(key string, args ...interface{}) string {
	return i.Locale(i.defaultLocale).Trans(key, args...)
}
