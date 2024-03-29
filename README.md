# Account Deposit Validator

In finance, it's common for accounts to have so-called "velocity limits". This program accepts or declines attempts to load funds into customers' accounts in real-time.

Each attempt to load funds will come as a single-line JSON payload, structured as follows:

```json
{
  "id": "1234",
  "customer_id": "1234",
  "load_amount": "$123.45",
  "time": "2018-01-01T00:00:00Z"
}
```

Each customer is subject to three limits:

- A maximum of $5,000 can be loaded per day
- A maximum of $20,000 can be loaded per week
- A maximum of 3 loads can be performed per day, regardless of amount

As such, a user attempting to load $3,000 twice in one day would be declined on the second attempt, as would a user attempting to load $400 four times in a day.

For each load attempt, a JSON response indicating whether the fund load was accepted based on the user's activity is returned, with the structure:

```json
{ "id": "1234", "customer_id": "1234", "accepted": true }
```

This project assumes the input arrives in ascending chronological order and that if a load ID is observed more than once for a particular user, all but the first instance is ignored. Each day is considered to end at midnight UTC, and weeks start on Monday (i.e. one second after 23:59:59 on Sunday).

## Implementation

The program reads `input.txt` and creates a Deposit struct from each line of JSON input. If the JSON is improperly formatted, or cannot be unmarshalled to a Deposit, then the program exits due to a fatal error. If the input is properly formatted then the program checks to see if the deposit has already been validated. If the load ID and customer ID have been validated previously, the input is skipped. Otherwise the deposit is validated and the response JSON is written to `output.txt`.

Deposits are validated with the help of daily and weekly ledgers. There is a daily and weekly ledger for each individual customer, and each ledger records the amount of money deposited into the customer's account during the time period. The daily ledger also records the total number of deposits for the day. Since the deposits are all received in chronological order the ledgers are reset whenever a customer makes a deposit during a new time period. If the deposit is valid the ledgers are updated to include the new deposit.

## Testing

Tests can be run along with a test coverage report by running `go test -v -cover`

## Execution

To run the program, clone the repository, compile the program using `go build` and run the executable.
