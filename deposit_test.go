package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsUnique_ShouldReturnTrueIfTheDepositHasNotBeenValidated(t *testing.T) {
	deposit := Deposit{"1", "1", "$1.00", time.Date(2021, 1, 9, 10, 0, 0, 0, time.UTC), 0}

	assert.True(t, deposit.IsUnique())
}

func TestValidate_ShouldReturnTrueIfDepositIsValid(t *testing.T) {
	first := Deposit{"1", "1", "$1.00", time.Date(2021, 1, 9, 10, 0, 0, 0, time.UTC), 0}

	assert.True(t, first.Validate())
}

func TestIsUnique_ShouldReturnFalseAfterTheDepositHasBeenValidated(t *testing.T) {
	deposit := Deposit{"1", "1", "$1.00", time.Date(2021, 1, 9, 10, 0, 0, 0, time.UTC), 0}

	assert.False(t, deposit.IsUnique())
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
