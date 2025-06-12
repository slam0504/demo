package i18n

import (
	"demo/internal/domain/card"
	"github.com/google/uuid"
	"testing"
)

type dummyCard struct{}

func TestTranslate(t *testing.T) {
	if Translate("zh", "invalid_id") != "無效的ID" {
		t.Fatalf("unexpected translation")
	}
	if Translate("fr", "invalid_id") != "invalid id" { // fallback to en
		t.Fatalf("fallback failed")
	}
	if Translate("unknown", "missing") != "missing" {
		t.Fatalf("missing key fallback failed")
	}
}

func TestTranslateCardNil(t *testing.T) {
	if TranslateCard("en", nil) != nil {
		t.Fatal("expected nil map")
	}
}

func TestTranslateCard(t *testing.T) {
	c := &card.Card{ID: uuid.New(), Name: "N", Cost: 1, Faction: "F", Category: "C", SubCategory: "S", Description: "D"}
	m := TranslateCard("zh", c)
	if m["名稱"] != "N" || m["費用"] != 1 {
		t.Fatalf("unexpected map %#v", m)
	}
	list := TranslateCards("zh", []*card.Card{c})
	if len(list) != 1 || list[0]["名稱"] != "N" {
		t.Fatalf("unexpected list %#v", list)
	}
}
