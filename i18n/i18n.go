package i18n

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var localizer = initLocal()

func initLocal() map[string]*i18n.Localizer {
	return map[string]*i18n.Localizer{}
}

func InitI18n(langs []string) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	for _, lang := range langs {
		_, err := bundle.LoadMessageFile(fmt.Sprintf("i18n/%v.json", lang))
		if err != nil {
			log.Fatal(err.Error())
		}
		localizer[lang] = i18n.NewLocalizer(bundle, lang)
	}
}

func Localize(message string, lang string) string {
	l := localizer[lang]
	if l == nil {
		return message
	}
	m, err := localizer[lang].Localize(&i18n.LocalizeConfig{MessageID: message})
	if err != nil {
		return message
	}
	return m
}

func LocalizeWithValue(message string, lang string, values map[string]string) string {
	l := localizer[lang]
	if l == nil {
		return message
	}
	m, err := localizer[lang].Localize(&i18n.LocalizeConfig{MessageID: message, TemplateData: values})
	if err != nil {
		return message
	}
	return m
}
