package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/travisbale/deposit-validator/deposit"
)

func TestProcessInput(t *testing.T) {
	validator := deposit.NewValidator()

	t.Run("processInput should return properly formatted JSON", func(t *testing.T) {
		input := `{"id":"15887","customer_id":"528","load_amount":"$3318.47","time":"2000-01-01T00:00:00Z"}`
		result, _ := processInput(validator, input)

		assert.Equal(t, result, `{"id":"15887","customer_id":"528","accepted":true}`)
	})

	t.Run("prcessInput should return an error if the deposit has been validated", func(t *testing.T) {
		input := `{"id":"15887","customer_id":"528","load_amount":"$3318.47","time":"2000-01-01T00:00:00Z"}`
		_, err := processInput(validator, input)

		assert.EqualError(t, err, "Deposit has already been processed")
	})
}
