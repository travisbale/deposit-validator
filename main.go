package main

import (
	"bufio"
	"encoding/json"
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
		var deposit Deposit

		// Parse the JSON payload into a Deposit
		err = json.Unmarshal([]byte(scanner.Text()), &deposit)
		checkError(err)

		// Ignore the deposit if it has already been validated
		if !deposit.HasBeenValidated() {
			// Create the response JSON and write it to the output file
			response := fmt.Sprintf(`{"id":"%s","customer_id":"%s","accepted":%t}`, deposit.ID, deposit.CustomerID, deposit.Validate())
			_, err := outFile.WriteString(response + "\n")
			checkError(err)
		}
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
