package main

import (
	"strconv"
	"strings"
	"time"
)

// Declare the velocity limits
const dailyLimit = 5000
const weeklyLimit = 20000
const maxDailyDeposits = 3

// A Deposit request is used to attempt to load funds into a customer account
type Deposit struct {
	ID           string    `json:"id"`
	CustomerID   string    `json:"customer_id"`
	Amount       string    `json:"load_amount"`
	Time         time.Time `json:"time"`
	ParsedAmount float64
}

// The dailyLedger is used to record and validate deposits on the current day
type dailyLedger struct {
	year     int
	month    time.Month
	day      int
	deposits int
	total    float64
}

// The weeklyLedger is used to record and validate deposits over the current week
type weeklyLedger struct {
	year  int
	week  int
	total float64
}

// Record all processed deposits to prevent duplicates
var validatedDeposits = make(map[string]bool)

// Keep daily and weekly ledgers for each customer
var dailyLedgers = make(map[string]dailyLedger)
var weeklyLedgers = make(map[string]weeklyLedger)

// IsUnique returns whether or not the deposit has already been processed
func (deposit *Deposit) IsUnique() bool {
	return !validatedDeposits[deposit.getUniqueIdentifier()]
}

// Validate returns whether or not the deposit is valid
func (deposit *Deposit) Validate() bool {
	err := deposit.parseAmount()

	// Record the deposit so it does not get processed twice
	validatedDeposits[deposit.getUniqueIdentifier()] = true

	if err == nil && deposit.validateDailyLimits() && deposit.validateWeeklyLimit() {
		// Record the deposit in the daily ledger
		dailyLedger := dailyLedgers[deposit.CustomerID]
		dailyLedger.deposits++
		dailyLedger.total += deposit.ParsedAmount
		dailyLedgers[deposit.CustomerID] = dailyLedger

		// Record the deposit in the weekly ledger
		weeklyLedger := weeklyLedgers[deposit.CustomerID]
		weeklyLedger.total += deposit.ParsedAmount
		weeklyLedgers[deposit.CustomerID] = weeklyLedger

		return true
	}

	return false
}

func (deposit *Deposit) getUniqueIdentifier() string {
	return deposit.ID + "-" + deposit.CustomerID
}

func (deposit *Deposit) validateDailyLimits() bool {
	// Get the customer's ledger, or an empty ledger if their ledger didn't exist
	ledger := dailyLedgers[deposit.CustomerID]

	year, month, day := deposit.Time.Date()

	// Update the current date and clear the ledger if the deposit occurs on a new day
	if ledger.year != year || ledger.month != month || ledger.day != day {
		ledger.year = year
		ledger.month = month
		ledger.day = day
		ledger.deposits = 0
		ledger.total = 0
		dailyLedgers[deposit.CustomerID] = ledger
	}

	return ledger.deposits < maxDailyDeposits && (ledger.total+deposit.ParsedAmount) <= dailyLimit
}

func (deposit *Deposit) validateWeeklyLimit() bool {
	// Get the customer's ledger, or an empty ledger if their ledger didn't exist
	ledger := weeklyLedgers[deposit.CustomerID]

	year, week := deposit.Time.ISOWeek()

	// Update the current week and clear the ledger if the deposit occurs in a new week
	if ledger.year != year || ledger.week != week {
		ledger.year = year
		ledger.week = week
		ledger.total = 0
		weeklyLedgers[deposit.CustomerID] = ledger
	}

	return (ledger.total + deposit.ParsedAmount) <= weeklyLimit
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
