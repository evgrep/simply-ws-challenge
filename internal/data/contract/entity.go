package contract

import "encoding/json"

// Entity - represents the data extracted from the database for a company entity
type Entity struct {
	Name      string      `json:"name"`
	Symbol    string      `json:"symbol"`
	Snowflake int32       `json:"snowflake"`
	Prices    []LastPrice `json:"last_prices,omitempty"`
}

// LastPrice - represents the JSON struct to return a last price from the database
type LastPrice struct {
	Date  string      `json:"date"`
	Price json.Number `json:"price"`
}
