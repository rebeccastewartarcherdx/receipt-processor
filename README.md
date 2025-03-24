# Receipt Processor

A simple http service that will process and calculate points for a receipt.

## How to Run
Clone the repo, then when in `receipt-processor`, run `go run .` in your terminal. This assumes you have go version 1.24 installed, this can be modified in `go.mod`. This will start the server at `localhost:6790`. Note this port can be modified in `main.go`.

To run tests, run `go test ./...` in your terminal.

## Endpoints

### Endpoint: Process Receipts

* Path: `{PORT}/receipts/process`
* Method: `POST`
* Payload: Receipt JSON
* Response: JSON containing an id for the receipt.

Description:

Takes in a JSON receipt and returns a JSON object with an ID.

The ID returned is the ID that can be passed into `{PORT}/receipts/{id}/points` to get the number of points the receipt was awarded.

How many points awarded are defined by the challenge spec.

Example Response:
```json
{ "id": "7fb1377b-b223-49d9-a31a-5a02701dd310" }
```

### Endpoint: Get Points
* Path: `/receipts/{id}/points`
* Method: `GET`
* Response: A JSON object containing the number of points awarded.

A simple Getter endpoint that looks up the receipt by the ID and returns an object specifying the points awarded.

Example Response:
```json
{ "points": 32 }
```