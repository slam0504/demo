package i18n

import (
	"embed"
	"encoding/json"
	"strings"

	"demo/internal/domain/card"
)

//go:embed en.json zh.json
var fs embed.FS

var translations map[string]map[string]string

func init() {
	translations = make(map[string]map[string]string)
	for _, lang := range []string{"en", "zh"} {
		data, err := fs.ReadFile(lang + ".json")
		if err != nil {
			continue
		}
		var m map[string]string
		if err := json.Unmarshal(data, &m); err == nil {
			translations[lang] = m
		}
	}
}

// Translate returns the localized message for the given key and language.
func Translate(lang, key string) string {
	lang = strings.ToLower(lang)
	if m, ok := translations[lang]; ok {
		if msg, ok := m[key]; ok {
			return msg
		}
	}
	// fallback to base language if region specific not found
	if idx := strings.Index(lang, "-"); idx != -1 {
		base := lang[:idx]
		if m, ok := translations[base]; ok {
			if msg, ok := m[key]; ok {
				return msg
			}
		}
	}
	if m, ok := translations["en"]; ok {
		if msg, ok := m[key]; ok {
			return msg
		}
	}
	return key
}

// TranslateCard returns a map with localized field names for the given card.
func TranslateCard(lang string, c *card.Card) map[string]interface{} {
	if c == nil {
		return nil
	}
	return map[string]interface{}{
		Translate(lang, "id"):          c.ID.String(),
		Translate(lang, "name"):        c.Name,
		Translate(lang, "cost"):        c.Cost,
		Translate(lang, "faction"):     c.Faction,
		Translate(lang, "category"):    c.Category,
		Translate(lang, "subcategory"): c.SubCategory,
		Translate(lang, "description"): c.Description,
	}
}

// TranslateCards localizes the field names for a slice of cards.
func TranslateCards(lang string, cards []*card.Card) []map[string]interface{} {
	res := make([]map[string]interface{}, 0, len(cards))
	for _, c := range cards {
		res = append(res, TranslateCard(lang, c))
	}
	return res
}
