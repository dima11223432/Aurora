// Package models defines domain types used across the API Gateway service.
package models

import "time"

// Stock represents a financial instrument mentioned in a post.
type Stock struct {
	StockName string `json:"stock_name"`
	Side      string `json:"side"`
}

// Post represents a social media post with associated stock mentions.
type Post struct {
	Stocks          []Stock   `json:"stocks"`
	PostText        string    `json:"post_text"`
	PostURI         string    `json:"post_uri"`
	ChannelUsername string    `json:"channel_username"`
	Reasoning       string    `json:"reasoning"`
	Date            time.Time `json:"date"`
}
