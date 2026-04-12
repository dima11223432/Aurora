package models

import "time"

type Stock struct {
	StockName string `json:"stock_name"`
	Side      string `json:"side"`
}

type AnalysedData struct {
	Stocks          []Stock   `json:"stocks"`
	PostText        string    `json:"post_text"`
	PostURI         string    `json:"post_uri"`
	ChannelUsername string    `json:"channel_username"`
	Date            time.Time `json:"date"`
	Reasoning       string    `json:"reasoning"`
}
