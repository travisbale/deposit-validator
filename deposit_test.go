package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type uniqueTest struct {
	deposit  Deposit
	expected bool
}

var uniqueTests = []uniqueTest{
	{Deposit{"1", "1", "$1.00", time.Date(2021, 1, 9, 10, 0, 0, 0, time.UTC), 0}, true},
	{Deposit{"1", "1", "$1.00", time.Date(2021, 1, 9, 10, 1, 0, 0, time.UTC), 0}, false},
	{Deposit{"1", "2", "$1.00", time.Date(2021, 1, 9, 10, 2, 0, 0, time.UTC), 0}, true},
	{Deposit{"2", "2", "$1.00", time.Date(2021, 1, 9, 10, 3, 0, 0, time.UTC), 0}, true},
	{Deposit{"3", "3", "$1.00", time.Date(2021, 1, 9, 10, 4, 0, 0, time.UTC), 0}, true},
	{Deposit{"1", "2", "$1.00", time.Date(2021, 1, 9, 10, 5, 0, 0, time.UTC), 0}, false},
}

func TestIsUnique(t *testing.T) {
	for _, test := range uniqueTests {
		if result := test.deposit.IsUnique(); result != test.expected {
			t.Errorf("Result '%t' does not match expected value '%t'", result, test.expected)
		}
	}
}

func TestValidate_ShouldReturnTrueIfDepositIsValid(t *testing.T) {
	first := Deposit{"1", "1", "$1.00", time.Date(2021, 1, 9, 10, 0, 0, 0, time.UTC), 0}

	assert.True(t, first.Validate())
}

func TestValidate_ShouldReturnTrueIfCustomerDepositsLessThanFourTimesADay(t *testing.T) {
	second := Deposit{"2", "1", "$1.00", time.Date(2021, 1, 9, 12, 1, 0, 0, time.UTC), 0}
	third := Deposit{"3", "1", "$1.00", time.Date(2021, 1, 9, 14, 0, 0, 0, time.UTC), 0}

	assert.True(t, second.Validate())
	assert.True(t, third.Validate())
}

func TestValidate_ShouldReturnFalseIfCustomerDepositsFourOrMoreTimesADay(t *testing.T) {
	fourth := Deposit{"4", "1", "$1.00", time.Date(2021, 1, 9, 23, 59, 59, 0, time.UTC), 0}
	fifth := Deposit{"5", "1", "$1.00", time.Date(2021, 1, 9, 23, 59, 59, 1, time.UTC), 0}

	assert.False(t, fourth.Validate())
	assert.False(t, fifth.Validate())
}

func TestValidate_CustomerDepositLimitIsResetDaily(t *testing.T) {
	deposit := Deposit{"6", "1", "$1.00", time.Date(2021, 1, 10, 0, 0, 0, 0, time.UTC), 0}

	assert.True(t, deposit.Validate())
}

func TestValidate_ShouldReturnTrueIfCustomerDepositsLessThanTheDailyLimitInADay(t *testing.T) {
	first := Deposit{"1", "2", "$4000.00", time.Date(2021, 1, 9, 10, 0, 0, 0, time.UTC), 0}
	second := Deposit{"2", "2", "$1000.00", time.Date(2021, 1, 9, 11, 0, 0, 0, time.UTC), 0}

	assert.True(t, first.Validate())
	assert.True(t, second.Validate())
}

func TestValidate_ShouldReturnFalseIfCustomerDepositsMoreThanTheDailyLimitInADay(t *testing.T) {
	third := Deposit{"3", "2", "$0.01", time.Date(2021, 1, 9, 23, 59, 59, 0, time.UTC), 0}

	assert.False(t, third.Validate())
}

func TestValidate_ShouldReturnTrueIfCustomerDepositsLessThanTheWeeklyLimitInAWeek(t *testing.T) {
	first := Deposit{"1", "3", "$5000.00", time.Date(2021, 1, 4, 0, 0, 0, 0, time.UTC), 0}
	second := Deposit{"2", "3", "$5000.00", time.Date(2021, 1, 5, 0, 0, 0, 0, time.UTC), 0}
	third := Deposit{"3", "3", "$5000.00", time.Date(2021, 1, 6, 0, 0, 0, 0, time.UTC), 0}
	fourth := Deposit{"4", "3", "$5000.00", time.Date(2021, 1, 7, 0, 0, 0, 0, time.UTC), 0}
	fifth := Deposit{"5", "3", "$0.01", time.Date(2021, 1, 10, 23, 59, 59, 0, time.UTC), 0}

	assert.True(t, first.Validate())
	assert.True(t, second.Validate())
	assert.True(t, third.Validate())
	assert.True(t, fourth.Validate())
	assert.False(t, fifth.Validate())
}

func TestValidate_ShouldReturnFalseIfAmountCannotBeParsed(t *testing.T) {
	deposit := Deposit{"1", "4", "%5000.00", time.Date(2021, 9, 0, 0, 0, 0, 0, time.UTC), 0}

	assert.False(t, deposit.Validate())
}
