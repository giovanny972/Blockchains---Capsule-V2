#!/bin/bash
# Script d'initialisation et démarrage automatique de Capsule Network
# Usage: ./start-blockchain.sh

set -e

CHAIN_ID="capsule-mainnet-1"
MONIKER="capsule-node"
HOME_DIR="/root/.simapp"

echo "🚀 Starting Capsule Network Blockchain..."
echo "=================================="

# Vérifier si déjà initialisé
if [ -f "$HOME_DIR/config/genesis.json" ]; then
    echo "✅ Blockchain already initialized"
    echo "📊 Starting existing blockchain..."
    exec simd start --home "$HOME_DIR"
fi

echo "🔧 First time setup - Initializing blockchain..."

# 1. Initialiser la blockchain
echo "⚙️  Step 1/5: Initializing node..."
simd init "$MONIKER" --chain-id "$CHAIN_ID" --home "$HOME_DIR"

# 2. Créer un compte validator
echo "🔑 Step 2/5: Creating validator account..."
simd keys add validator --keyring-backend test --home "$HOME_DIR" 2>&1 | tee /tmp/validator_key.txt

VALIDATOR_ADDRESS=$(simd keys show validator -a --keyring-backend test --home "$HOME_DIR")
echo "Validator address: $VALIDATOR_ADDRESS"

# 3. Ajouter le compte au genesis avec des fonds
echo "💰 Step 3/5: Adding genesis account..."
# Total Supply: 1B STAKE + 100M UCAPS
# Validator allocation: 500M STAKE (50% for staking + reserve) + 10M UCAPS (10% ecosystem)
simd genesis add-genesis-account "$VALIDATOR_ADDRESS" 500000000stake,10000000ucaps --home "$HOME_DIR"

# 4. Créer la transaction genesis pour le validator
echo "🏛️  Step 4/5: Creating genesis validator..."
# Stake 200M (40% of validator allocation) to secure the network
simd genesis gentx validator 200000000stake \
    --chain-id "$CHAIN_ID" \
    --keyring-backend test \
    --home "$HOME_DIR"

# 5. Collecter les gentx
echo "📝 Step 5/5: Collecting genesis transactions..."
simd genesis collect-gentxs --home "$HOME_DIR"

# Configuration optionnelle pour faciliter le développement
echo "⚙️  Configuring node settings..."

# Activer API REST
sed -i 's/enable = false/enable = true/g' "$HOME_DIR/config/app.toml"

# Activer Swagger
sed -i 's/swagger = false/swagger = true/g' "$HOME_DIR/config/app.toml"

# Configuration RPC pour accepter les connexions externes
sed -i 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g' "$HOME_DIR/config/config.toml"

# Configuration API pour accepter les connexions externes
sed -i 's/address = "tcp:\/\/localhost:1317"/address = "tcp:\/\/0.0.0.0:1317"/g' "$HOME_DIR/config/app.toml"

# Activer CORS
sed -i 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/g' "$HOME_DIR/config/app.toml"

echo ""
echo "✅ Blockchain initialized successfully!"
echo "=================================="
echo "Chain ID: $CHAIN_ID"
echo "Moniker: $MONIKER"
echo "Validator Address: $VALIDATOR_ADDRESS"
echo ""
echo "📊 Starting blockchain..."
echo "=================================="

# Démarrer la blockchain
exec simd start --home "$HOME_DIR"
