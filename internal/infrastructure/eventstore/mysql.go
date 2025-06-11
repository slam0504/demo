package eventstore

import (
	"context"
	"encoding/json"
	"fmt"

	"demo/internal/domain/card"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type EventRecord struct {
	ID      uint   `gorm:"primaryKey"`
	CardID  string `gorm:"index"`
	Type    string
	Payload []byte
}

type MySQLStore struct {
	DB *gorm.DB
}

// NewMySQLStore creates a GORM-based MySQL event store.
func NewMySQLStore(dsn string) (*MySQLStore, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&EventRecord{}); err != nil {
		return nil, err
	}
	return &MySQLStore{DB: db}, nil
}

func eventCardID(evt interface{}) (string, error) {
	switch e := evt.(type) {
	case card.CardCreated:
		return e.ID.String(), nil
	case card.CardUpdated:
		return e.ID.String(), nil
	default:
		return "", fmt.Errorf("unknown event type %T", evt)
	}
}

// Save stores events in MySQL.
func (s *MySQLStore) Save(ctx context.Context, events []interface{}) error {
	if len(events) == 0 {
		return nil
	}
	for _, evt := range events {
		id, err := eventCardID(evt)
		if err != nil {
			return err
		}
		data, err := json.Marshal(evt)
		if err != nil {
			return err
		}
		rec := EventRecord{CardID: id, Type: fmt.Sprintf("%T", evt), Payload: data}
		if err := s.DB.WithContext(ctx).Create(&rec).Error; err != nil {
			return err
		}
	}
	return nil
}

// Load rebuilds the card state from events.
func (s *MySQLStore) Load(ctx context.Context, id string) (*card.Card, error) {
	var records []EventRecord
	if err := s.DB.WithContext(ctx).Where("card_id = ?", id).Order("id").Find(&records).Error; err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, nil
	}
	c := &card.Card{}
	for _, r := range records {
		switch r.Type {
		case "card.CardCreated":
			var evt card.CardCreated
			if err := json.Unmarshal(r.Payload, &evt); err != nil {
				return nil, err
			}
			c.ID = evt.ID
			c.Name = evt.Name
			c.Cost = evt.Cost
			c.Faction = evt.Faction
			c.Category = evt.Category
			c.SubCategory = evt.SubCategory
			c.Description = evt.Description
		case "card.CardUpdated":
			var evt card.CardUpdated
			if err := json.Unmarshal(r.Payload, &evt); err != nil {
				return nil, err
			}
			c.Name = evt.Name
			c.Cost = evt.Cost
			c.Faction = evt.Faction
			c.Category = evt.Category
			c.SubCategory = evt.SubCategory
			c.Description = evt.Description
		}
	}
	return c, nil
}

// Search loads all cards and filters them.
func (s *MySQLStore) Search(ctx context.Context, name string, cost int, faction, category, sub string) ([]*card.Card, error) {
	var ids []string
	if err := s.DB.WithContext(ctx).Model(&EventRecord{}).Distinct("card_id").Find(&ids).Error; err != nil {
		return nil, err
	}
	var cards []*card.Card
	for _, id := range ids {
		c, err := s.Load(ctx, id)
		if err != nil || c == nil {
			continue
		}
		if (name == "" || c.Name == name) &&
			(cost == 0 || c.Cost == cost) &&
			(faction == "" || c.Faction == faction) &&
			(category == "" || c.Category == category) &&
			(sub == "" || c.SubCategory == sub) {
			tmp := *c
			cards = append(cards, &tmp)
		}
	}
	return cards, nil
}
