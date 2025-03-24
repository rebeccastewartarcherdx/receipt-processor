package models

import "time"

type Receipt struct {
	Retailer         string
	PurchaseDateTime time.Time
	Total            int
	Items            []Item
	Points           int
}

type Item struct {
	ShortDescription string
	Price            int
}

type UnprocessedReceipt struct {
	Retailer     string            `json:"retailer"`
	PurchaseDate string            `json:"purchaseDate"`
	PurchaseTime string            `json:"purchaseTime"`
	Total        string            `json:"total"`
	Items        []UnprocessedItem `json:"items"`
}

type UnprocessedItem struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}
