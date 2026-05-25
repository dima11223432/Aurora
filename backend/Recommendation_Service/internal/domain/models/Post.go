// Package models defines domain types for the Recommendation Service.
package models

import "time"

// Stock represents a financial instrument mentioned in a post.
type Stock struct {
	StockName string `json:"stock_name"`
	Side      string `json:"side"`
}

// Post represents a recommended social media post with stock mentions.
type Post struct {
	Stocks          []Stock   `json:"stocks"`
	PostText        string    `json:"post_text"`
	PostURI         string    `json:"post_uri"`
	ChannelUsername string    `json:"channel_username"`
	Reasoning       string    `json:"reasoning"`
	Date            time.Time `json:"date"`
}
