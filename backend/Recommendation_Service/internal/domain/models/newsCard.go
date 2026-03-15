package models

import "time"

type Stock struct {
	StockName string `json:"stock_name"`
	Side      string `json:"side"`
}

type NewsCard struct {
	Stocks   []Stock   `json:"stocks"`
	PostText string    `json:"post_text"`
	PostURI  string    `json:"post_uri"`
	Date     time.Time `json:"date"`
}
