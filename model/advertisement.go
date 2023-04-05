package model

import "time"

type Advertisement struct {
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Author      string    `json:"author"`
	Date        time.Time `json:"date"`
}
