package i18n

import ()

type Tupel struct {
	Key   string
	Value string
}

type Locale struct {
	Lang   string
	Tupels []Tupel
}

var (
	Locales []Locale
)

func LoadLocales() {
}
