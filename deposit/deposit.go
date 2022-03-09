package deposit

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

// A Deposit request is used to attempt to load funds into a customer account
type Deposit struct {
	ID           string    `json:"id"`
	CustomerID   string    `json:"customer_id"`
	Amount       string    `json:"load_amount"`
	Time         time.Time `json:"time"`
	ParsedAmount float64
}

func ParseJson(depositJson string) (*Deposit, error) {
	var deposit Deposit

	if err := json.Unmarshal([]byte(depositJson), &deposit); err != nil {
		return nil, err
	}

	if err := deposit.parseAmount(); err != nil {
		return nil, err
	}

	return &deposit, nil
}

func (deposit *Deposit) parseAmount() error {
	amount, err := strconv.ParseFloat(strings.TrimPrefix(deposit.Amount, "$"), 64)
	if err != nil {
		// Input wasn't properly formatted
		return err
	}

	// Save the amount to the struct so it only has to be calculated once
	deposit.ParsedAmount = amount

	return nil
}
