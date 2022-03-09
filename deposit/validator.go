package deposit

import "time"

// Declare the velocity limits
const dailyLimit = 5000
const weeklyLimit = 20000
const maxDailyDeposits = 3

type Validator interface {
	HasBeenValidated(deposit *Deposit) bool
	Validate(deposit *Deposit) bool
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

type validator struct {
	// Record all validated deposits to prevent duplicates
	validatedDeposits map[string]bool

	// Keep daily and weekly ledgers for each customer
	dailyLedgers  map[string]dailyLedger
	weeklyLedgers map[string]weeklyLedger
}

func NewValidator() Validator {
	return &validator{
		validatedDeposits: make(map[string]bool),
		dailyLedgers:      make(map[string]dailyLedger),
		weeklyLedgers:     make(map[string]weeklyLedger),
	}
}

// HasBeenValidated returns whether or not the deposit has already been processed
func (v *validator) HasBeenValidated(deposit *Deposit) bool {
	return v.validatedDeposits[getUniqueIdentifier(deposit)]
}

// Validate returns whether or not the deposit is valid
func (v *validator) Validate(deposit *Deposit) bool {
	err := deposit.parseAmount()

	// Record the deposit so it does not get processed twice
	v.validatedDeposits[getUniqueIdentifier(deposit)] = true

	if err == nil && v.validateDailyLimits(deposit) && v.validateWeeklyLimit(deposit) {
		// Record the deposit in the daily ledger
		dailyLedger := v.dailyLedgers[deposit.CustomerID]
		dailyLedger.deposits++
		dailyLedger.total += deposit.ParsedAmount
		v.dailyLedgers[deposit.CustomerID] = dailyLedger

		// Record the deposit in the weekly ledger
		weeklyLedger := v.weeklyLedgers[deposit.CustomerID]
		weeklyLedger.total += deposit.ParsedAmount
		v.weeklyLedgers[deposit.CustomerID] = weeklyLedger

		return true
	}

	return false
}

func getUniqueIdentifier(deposit *Deposit) string {
	return deposit.ID + "-" + deposit.CustomerID
}

func (v *validator) validateDailyLimits(deposit *Deposit) bool {
	// Get the customer's ledger, or an empty ledger if their ledger didn't exist
	ledger := v.dailyLedgers[deposit.CustomerID]

	year, month, day := deposit.Time.Date()

	// Update the current date and clear the ledger if the deposit occurs on a new day
	if ledger.year != year || ledger.month != month || ledger.day != day {
		ledger.year = year
		ledger.month = month
		ledger.day = day
		ledger.deposits = 0
		ledger.total = 0
		v.dailyLedgers[deposit.CustomerID] = ledger
	}

	return ledger.deposits < maxDailyDeposits && (ledger.total+deposit.ParsedAmount) <= dailyLimit
}

func (v *validator) validateWeeklyLimit(deposit *Deposit) bool {
	// Get the customer's ledger, or an empty ledger if their ledger didn't exist
	ledger := v.weeklyLedgers[deposit.CustomerID]

	year, week := deposit.Time.ISOWeek()

	// Update the current week and clear the ledger if the deposit occurs in a new week
	if ledger.year != year || ledger.week != week {
		ledger.year = year
		ledger.week = week
		ledger.total = 0
		v.weeklyLedgers[deposit.CustomerID] = ledger
	}

	return (ledger.total + deposit.ParsedAmount) <= weeklyLimit
}
