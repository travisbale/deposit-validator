package deposit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var v Validator

func setupTestCase(t *testing.T) func(t *testing.T) {
	// Reset the ledgers and list of validated deposits between tests
	v = NewValidator()

	// Return a function to perform test teardown
	return func(t *testing.T) {
		// Nothing to tear down for these tests
	}
}

func TestHasBeenValidated(t *testing.T) {
	tearDownTestCase := setupTestCase(t)
	defer tearDownTestCase(t)

	deposit := Deposit{"1", "1", "$1.00", time.Date(2021, 1, 9, 10, 0, 0, 0, time.UTC), 0}

	t.Run("HasBeenValidated should return false if the deposit has not been validated", func(t *testing.T) {
		assert.False(t, v.HasBeenValidated(&deposit))
	})

	t.Run("HasBeenValidated should return true if the deposit has been validated", func(t *testing.T) {
		v.Validate(&deposit)
		assert.True(t, v.HasBeenValidated(&deposit))
	})
}

func TestValidate_NumberOfDailyDeposits(t *testing.T) {
	tearDownTestCase := setupTestCase(t)
	defer tearDownTestCase(t)

	t.Run("Validate should return true if customer deposits less than four times in a day", func(t *testing.T) {
		first := Deposit{"1", "1", "$1.00", time.Date(2021, 1, 9, 10, 0, 0, 0, time.UTC), 0}
		second := Deposit{"2", "1", "$1.00", time.Date(2021, 1, 9, 12, 1, 0, 0, time.UTC), 0}
		third := Deposit{"3", "1", "$1.00", time.Date(2021, 1, 9, 14, 0, 0, 0, time.UTC), 0}
		assert.True(t, v.Validate(&first))
		assert.True(t, v.Validate(&second))
		assert.True(t, v.Validate(&third))
	})

	t.Run("Validate should return false if customer deposits four or more times in a day", func(t *testing.T) {
		fourth := Deposit{"4", "1", "$1.00", time.Date(2021, 1, 9, 23, 59, 59, 0, time.UTC), 0}
		fifth := Deposit{"5", "1", "$1.00", time.Date(2021, 1, 9, 23, 59, 59, 1, time.UTC), 0}
		assert.False(t, v.Validate(&fourth))
		assert.False(t, v.Validate(&fifth))
	})

	t.Run("Validate should return true once the ledger is reset the next day", func(t *testing.T) {
		sixth := Deposit{"6", "1", "$1.00", time.Date(2021, 1, 10, 0, 0, 0, 0, time.UTC), 0}
		assert.True(t, v.Validate(&sixth))
	})
}

func TestValidate_TotalValueOfDailyDeposits(t *testing.T) {
	tearDownTestCase := setupTestCase(t)
	defer tearDownTestCase(t)

	t.Run("Validate should return true if customer deposits less than the daily limit in a day", func(t *testing.T) {
		first := Deposit{"1", "2", "$4000.00", time.Date(2021, 1, 9, 10, 0, 0, 0, time.UTC), 0}
		second := Deposit{"2", "2", "$1000.00", time.Date(2021, 1, 9, 11, 0, 0, 0, time.UTC), 0}
		assert.True(t, v.Validate(&first))
		assert.True(t, v.Validate(&second))
	})

	t.Run("Validate should return false if the customer deposits more than the daily limit in a day", func(t *testing.T) {
		third := Deposit{"3", "2", "$0.01", time.Date(2021, 1, 9, 23, 59, 59, 0, time.UTC), 0}
		assert.False(t, v.Validate(&third))
	})

	t.Run("Validate should return true once the ledger is resset the next day", func(t *testing.T) {
		fourth := Deposit{"4", "2", "$0.01", time.Date(2021, 1, 10, 0, 0, 0, 0, time.UTC), 0}
		assert.True(t, v.Validate(&fourth))
	})
}

func TestValidate_TotalValueOfWeeklyDeposits(t *testing.T) {
	tearDownTestCase := setupTestCase(t)
	defer tearDownTestCase(t)

	t.Run("Validate should return true if customer deposits less than the weekly limit in a week", func(t *testing.T) {
		first := Deposit{"1", "3", "$5000.00", time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC), 0}
		second := Deposit{"2", "3", "$5000.00", time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC), 0}
		third := Deposit{"3", "3", "$5000.00", time.Date(2021, 1, 6, 0, 0, 0, 0, time.UTC), 0}
		fourth := Deposit{"4", "3", "$5000.00", time.Date(2021, 1, 7, 0, 0, 0, 0, time.UTC), 0}
		assert.True(t, v.Validate(&first))
		assert.True(t, v.Validate(&second))
		assert.True(t, v.Validate(&third))
		assert.True(t, v.Validate(&fourth))
	})

	t.Run("Validate should return false if the customer deposits more than the weekly limit in a week", func(t *testing.T) {
		fifth := Deposit{"5", "3", "$0.01", time.Date(2021, 1, 10, 23, 59, 59, 0, time.UTC), 0}
		assert.False(t, v.Validate(&fifth))
	})

	t.Run("Validate should return true once the ledger is reset the next day", func(t *testing.T) {
		sixth := Deposit{"6", "3", "$0.01", time.Date(2021, 1, 11, 0, 0, 0, 0, time.UTC), 0}
		assert.True(t, v.Validate(&sixth))
	})
}

func TestValidate_InvalidInput(t *testing.T) {
	tearDownTestCase := setupTestCase(t)
	defer tearDownTestCase(t)

	t.Run("Validate should return false if the amount cannot be parsed", func(t *testing.T) {
		deposit := Deposit{"1", "4", "%5000.00", time.Date(2021, 9, 0, 0, 0, 0, 0, time.UTC), 0}
		assert.False(t, v.Validate(&deposit))
	})
}
