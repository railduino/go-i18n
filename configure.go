package i18n

type Configuration struct {
	StaticPrefix string
	LocaleDir    string
}

var (
	Config = Configuration{
		StaticPrefix: "/static/",
		LocaleDir:    "./locales",
	}
)

func Configure(prefix, locales string) {
	Config.StaticPrefix = prefix
	Config.LocaleDir = locales
}
