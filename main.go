package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
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

	// Scan the input file line by line
	for scanner.Scan() {
		response, err := processInput(scanner.Text())

		if err == nil {
			_, err := outFile.WriteString(response + "\n")
			checkError(err)
		}
	}
}

func processInput(input string) (string, error) {
	var deposit Deposit

	// Parse the JSON input into a deposit
	err := json.Unmarshal([]byte(input), &deposit)
	checkError(err)

	// Ignore the deposit if it has already been validated
	if !deposit.HasBeenValidated() {
		result := fmt.Sprintf(`{"id":"%s","customer_id":"%s","accepted":%t}`, deposit.ID, deposit.CustomerID, deposit.Validate())
		return result, nil
	}

	return "", errors.New("Deposit has already been processed")
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
