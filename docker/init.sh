#!/bin/bash

# Configuration
CHAIN_ID=${CHAIN_ID:-"capsule-testnet-1"}
MONIKER=${MONIKER:-"capsule-validator"}
KEYRING_BACKEND=${KEYRING_BACKEND:-"test"}
VALIDATOR_KEY=${VALIDATOR_KEY:-"validator"}
USER_KEY=${USER_KEY:-"alice"}

echo "Initializing Capsule blockchain node..."
echo "Chain ID: $CHAIN_ID"
echo "Moniker: $MONIKER"

# Initialize the blockchain
simd init $MONIKER --chain-id $CHAIN_ID

# Create validator key
echo "Creating validator key..."
simd keys add $VALIDATOR_KEY --keyring-backend $KEYRING_BACKEND

# Create user key for testing
echo "Creating test user key..."
simd keys add $USER_KEY --keyring-backend $KEYRING_BACKEND

# Get addresses
VALIDATOR_ADDRESS=$(simd keys show $VALIDATOR_KEY -a --keyring-backend $KEYRING_BACKEND)
USER_ADDRESS=$(simd keys show $USER_KEY -a --keyring-backend $KEYRING_BACKEND)

echo "Validator address: $VALIDATOR_ADDRESS"
echo "User address: $USER_ADDRESS"

# Add accounts to genesis
echo "Adding accounts to genesis..."
simd genesis add-genesis-account $VALIDATOR_ADDRESS 100000000000000stake
simd genesis add-genesis-account $USER_ADDRESS 10000000000stake

# Create genesis transaction
echo "Creating genesis transaction..."
simd genesis gentx $VALIDATOR_KEY 100000000stake --keyring-backend $KEYRING_BACKEND --chain-id $CHAIN_ID

# Collect genesis transactions
echo "Collecting genesis transactions..."
simd genesis collect-gentxs

# Update configuration for Docker environment
echo "Updating configuration..."

# Enable API server
sed -i 's/enable = false/enable = true/' ~/.simd/config/app.toml
sed -i 's/swagger = false/swagger = true/' ~/.simd/config/app.toml

# Configure RPC server
sed -i 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/' ~/.simd/config/config.toml

# Configure gRPC server
sed -i 's/address = "127.0.0.1:9090"/address = "0.0.0.0:9090"/' ~/.simd/config/app.toml

# Configure REST API server
sed -i 's/address = "tcp:\/\/localhost:1317"/address = "tcp:\/\/0.0.0.0:1317"/' ~/.simd/config/app.toml

# Allow all CORS origins for development
sed -i 's/cors_allowed_origins = \[\]/cors_allowed_origins = \["\*"\]/' ~/.simd/config/config.toml

# Set minimum gas prices
sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "0.0001stake"/' ~/.simd/config/app.toml

# Enable unsafe CORS for development (remove in production)
sed -i 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/' ~/.simd/config/app.toml

echo "Initialization complete!"
echo "Validator key: $VALIDATOR_KEY"
echo "User key: $USER_KEY"
echo "Chain ID: $CHAIN_ID"

# Display key information
echo ""
echo "=== Key Information ==="
simd keys list --keyring-backend $KEYRING_BACKEND

echo ""
echo "=== Genesis File Info ==="
cat ~/.simd/config/genesis.json | jq '.chain_id, .genesis_time'

echo ""
echo "Node initialization completed successfully!"