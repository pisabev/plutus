# Plutus demo

### Setup environment variables

```bash
cp env.example .env
```
### Start the application

```bash
docker-compose up
```

### Run tests

```bash
go test ./...
```

### Write transactions

```bash
curl -X POST "http://localhost:8080/webhooks/transaction" -d '{"transactionId":"tqZi6QapS41zcEHy", "orderId":"c66oxMaisTwJQXjD", "transactionType":"SALE", "amount": "20.00", "currency":"EUR", "description":"Test transaction", "accountId":"001"}'
curl -X POST "http://localhost:8080/webhooks/transaction" -d '{"transactionId":"tqZi6QapS41zcEHy2", "orderId":"c66oxMaisTwJQXjD", "transactionType":"SALE", "amount": "20.00", "currency":"EUR", "description":"Test transaction", "accountId":"001"}'
curl -X POST "http://localhost:8080/webhooks/transaction" -d '{"transactionId":"tqZi6QapS41zcEHy3", "orderId":"c66oxMaisTwJQXjD", "transactionType":"CREDIT", "amount": "10.00", "currency":"EUR", "description":"Test transaction", "accountId":"001"}'
```

### Account balance

```bash
curl -X GET "http://localhost:8080/account/001"
```

{"balance":"30.00"}

### Some stress test with apache ab

```bash
ab -c 100 -n 2000 -p test.json -T application/json http://127.0.0.1:8080/webhooks/transaction
```

### Fetch all

```bash
curl -X GET "http://localhost:8080/all"
```