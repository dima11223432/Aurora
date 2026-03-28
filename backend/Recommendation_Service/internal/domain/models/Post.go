package models

import "time"

type Stock struct {
	StockName string `json:"stock_name"`
	Side      string `json:"side"`
}

type Post struct {
	Stocks          []Stock   `json:"stocks"`
	PostText        string    `json:"post_text"`
	PostURI         string    `json:"post_uri"`
	ChannelUsername string    `json:"channel_username"`
	Reasoning       string    `json:"reasoning"`
	Date            time.Time `json:"date"`
}
