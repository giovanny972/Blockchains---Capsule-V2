# Docker Deployment for Capsule Blockchain

## Quick Start

### Build and Run

```bash
# Build and start the node
docker-compose up -d

# Check logs
docker-compose logs -f capsule-node

# Stop the node
docker-compose down
```

### With Monitoring

```bash
# Start with Prometheus monitoring
docker-compose --profile monitoring up -d

# Access Prometheus at http://localhost:9091
```

## Manual Commands

### Build only
```bash
docker build -t capsule-blockchain .
```

### Run manually
```bash
docker run -d \
  --name capsule-node \
  -p 26656:26656 \
  -p 26657:26657 \
  -p 1317:1317 \
  -p 9090:9090 \
  -v capsule_data:/root/.simd \
  capsule-blockchain
```

## Endpoints

- **RPC**: http://localhost:26657
- **REST API**: http://localhost:1317
- **gRPC**: localhost:9090
- **P2P**: localhost:26656
- **Prometheus** (optional): http://localhost:9091

## Testing the Node

```bash
# Check node status
curl http://localhost:26657/status

# List accounts
docker exec capsule-blockchain simd keys list --keyring-backend test

# Check balances
docker exec capsule-blockchain simd query bank balances $(docker exec capsule-blockchain simd keys show alice -a --keyring-backend test)

# Send transaction
docker exec capsule-blockchain simd tx bank send alice $(docker exec capsule-blockchain simd keys show validator -a --keyring-backend test) 1000000stake --keyring-backend test --chain-id capsule-testnet-1 --yes

# Query timecapsule module
docker exec capsule-blockchain simd query timecapsule params
```

## Persistence

Data is persisted in the `capsule_data` Docker volume. To reset:

```bash
docker-compose down
docker volume rm capsule_capsule_data
docker-compose up -d
```

## Production Notes

- Change `KEYRING_BACKEND` from `test` to `file` or `os`
- Set proper minimum gas prices
- Configure firewall rules
- Use proper TLS certificates
- Disable unsafe CORS
- Set up proper backup for the data volume