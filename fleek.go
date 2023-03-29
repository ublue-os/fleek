package fleek

import (
	"embed"
	"os"

	"github.com/fitv/go-i18n"
)

//go:embed locales/*.yml
var fs embed.FS

func NewApp() *App {
	i18n, err := i18n.New(fs, "locales")
	if err != nil {
		panic(err)
	}
	i18n.SetDefaultLocale(locale())
	return &App{
		I18n: i18n,
	}
}

type App struct {
	*i18n.I18n
}

// locale returns the two digit locale code
// from the LANG environment variable, or "en"
// if unset.
func locale() string {
	lang := os.Getenv("LANG")
	if lang == "" {
		lang = "en"
	}
	locale := lang[:2]
	return locale
}
