# Blockchain-in-go

A simple blockchain implementation written in Go, featuring proof-of-work consensus mechanism and a REST API.

## Features

- Basic blockchain data structure
- Proof-of-Work consensus algorithm
- REST API for interacting with the blockchain
- Thread-safe blockchain operations

## API Endpoints

- `GET /blockchain` - Get the current state of the blockchain
- `POST /mine` - Mine a new block with provided data

## How to Run

1. Make sure you have Go installed (version 1.16 or higher)
2. Clone this repository
3. Run the application:

```
go run Blockchain.go
```

4. Access the blockchain via HTTP endpoints:

```
# View the blockchain
curl http://localhost:8080/blockchain

# Add a new block
curl -X POST -H "Content-Type: application/json" -d '{"data":"Your block data here"}' http://localhost:8080/mine
```

## Future Improvements

- Add persistent storage
- Implement peer-to-peer networking
- Add wallet and transaction functionality
- Dynamic difficulty adjustment
