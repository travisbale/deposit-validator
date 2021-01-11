# KOHO Take Home Assignment

## Instructions

In finance, it's common for accounts to have so-called "velocity limits". In this task, you'll write a program that accepts or declines attempts to load funds into customers' accounts in real-time.

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

For each load attempt, you should return a JSON response indicating whether the fund load was accepted based on the user's activity, with the structure:

```json
{ "id": "1234", "customer_id": "1234", "accepted": true }
```

You can assume that the input arrives in ascending chronological order and that if a load ID is observed more than once for a particular user, all but the first instance can be ignored. Each day is considered to end at midnight UTC, and weeks start on Monday (i.e. one second after 23:59:59 on Sunday).

Your program should process lines from `input.txt` and return output in the format specified above, either to standard output or a file. Expected output given our input data can be found in `output.txt`.

You're welcome to write your program in a general-purpose language of your choosing, but as we use Go on the back-end and TypeScript on the front-end, we do have a preference towards solutions written in Go (back-end) and TypeScript (front-end).

We value well-structured, self-documenting code with sensible test coverage. Descriptive function and variable names are appreciated, as is isolating your business logic from the rest of your code.

## Solution

I wasn't sure what the larger context of this problem was so I decided to go with a relatively straightforward implementation.

*Disclaimer:* This is my first foray into the world of Golang

### Implementation

The program reads `input.txt`, creating a Deposit structure from each line of JSON input. If the JSON is improperly formatted, or cannot be unmarshalled to a Deposit, then the program exits due to a fatal error. If the input is properly formatted, then the deposit is checked for uniqueness. If the load ID and customer ID have been processed previously, the input is skipped. If the deposit is unique, the deposit is validated and the response JSON is written to `output.txt`.

Deposits are validated with the help of daily and weekly "ledgers". There is a ledger for each individual customer, and each ledger records the amount of money deposited into the customer's account during the time period. The daily ledger also records the total number of deposits. Since the deposits are all received in chronological order they are reset whenever a deposit occurs during a new time period. If the deposit is valid, then the ledgers are updated to reflect the new deposit.

