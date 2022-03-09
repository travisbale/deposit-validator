package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/travisbale/deposit-validator/deposit"
)

func main() {
	// Open the input file for reading
	inFile, err := os.Open("input.txt")
	checkError(err)
	defer inFile.Close()

	// Open the output file for writing
	outFile, err := os.Create("output.txt")
	checkError(err)
	defer outFile.Close()

	scanner := bufio.NewScanner(inFile)
	depositValidator := deposit.NewValidator()

	// Scan the input file line by line
	for scanner.Scan() {
		response, err := processInput(depositValidator, scanner.Text())

		if err == nil {
			_, err := outFile.WriteString(response + "\n")
			checkError(err)
		}
	}
}

func processInput(depositValidator deposit.Validator, input string) (string, error) {
	deposit, err := deposit.ParseJson(input)
	checkError(err)

	// Ignore the deposit if it has already been validated
	if !depositValidator.HasBeenValidated(deposit) {
		isValid := depositValidator.Validate(deposit)
		result := fmt.Sprintf(`{"id":"%s","customer_id":"%s","accepted":%t}`, deposit.ID, deposit.CustomerID, isValid)
		return result, nil
	}

	return "", errors.New("deposit has already been processed")
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
