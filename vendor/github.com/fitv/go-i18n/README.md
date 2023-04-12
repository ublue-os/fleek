## I18n for Go

## Install
```
go get -u github.com/fitv/go-i18n
```

## Usage
YAML files
```
├── locales
│   ├── en.yml
│   ├── zh.yml
└── main.go
```

```go
package main

import "github.com/fitv/go-i18n"

//go:embed locales/*.yml
var fs embed.FS

func main() {
    i18n, err := i18n.New(fs, "locales")
    i18n.SetDefaultLocale("en")

    i18n.Trans("hello.world") // World
    i18n.Locale("zh").Trans("hello.world") // 世界

    // with params
    user := map[string]interface{}{"name": "Jack", "email": "jack@example.com"}
    // foo %s
    i18n.Trans("hello.foo", "bar") // foo bar
    // Name: {name}, Email: {email}
    i18n.Trans("user.description", user) // Name: Jack, Email: jack@example.com
}
```
