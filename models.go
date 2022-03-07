package main

import "time"

type Log struct {
	Type      string    `json:"type"`
	Message   string    `json:"message"`
	Level     string    `json:"level"`
	Event     string    `json:"event"`
	CreatedAt time.Time `json:"createdat"`
	UserID    uint      `json:"userid"`
	DeckID    uint      `json:"deckid"`
	CardID    uint      `json:"cardid"`
}
