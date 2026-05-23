// Package models defines domain types for the Cache Service.
package models

import "time"

// Stock represents a financial instrument mention.
type Stock struct {
	StockName string `json:"stock_name"`
	Side      string `json:"side"`
}

// AnalysedData represents analysed financial data from a social media post.
type AnalysedData struct {
	Stocks          []Stock   `json:"stocks"`
	PostText        string    `json:"post_text"`
	PostURI         string    `json:"post_uri"`
	ChannelUsername string    `json:"channel_username"`
	Date            time.Time `json:"date"`
	Reasoning       string    `json:"reasoning"`
}
