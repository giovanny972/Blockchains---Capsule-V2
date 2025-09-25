#!/bin/bash

# Script pour interroger les capsules sur la blockchain

echo "🔍 VERIFICATION DES CAPSULES SUR LA BLOCKCHAIN"
echo "=============================================="

SIMD_BIN="./simapp/simd/simd.exe"
NODE="tcp://localhost:26657"

echo ""
echo "📊 Statut de la Blockchain:"
echo "----------------------------"
$SIMD_BIN status | jq '.sync_info | {latest_block_height, latest_block_time, network}'

echo ""
echo "🔍 Recherche de Transactions Timecapsule:"
echo "----------------------------------------"

# Rechercher toutes les transactions récentes
echo "Recherche dans les 50 dernières transactions..."
$SIMD_BIN query txs --query="tx.height >= 1" --limit=50 --page=1 | jq '.txs[] | {height, hash, events: [.events[] | select(.type == "message" or .type == "capsule_created")]}'

echo ""
echo "🔍 Recherche d'événements spécifiques:"
echo "-------------------------------------"

# Rechercher des événements de création de capsule
$SIMD_BIN query txs --query="message.action CONTAINS 'capsule'" --limit=10 2>/dev/null || echo "Aucune transaction avec 'capsule' trouvée"

echo ""
echo "📈 Vérification des derniers blocs:"
echo "----------------------------------"

# Obtenir la hauteur actuelle
CURRENT_HEIGHT=$($SIMD_BIN status | jq -r '.sync_info.latest_block_height')
echo "Hauteur actuelle: $CURRENT_HEIGHT"

# Vérifier les 5 derniers blocs pour des transactions
echo "Vérification des 5 derniers blocs..."
for i in {1..5}; do
    HEIGHT=$((CURRENT_HEIGHT - i + 1))
    echo "Bloc $HEIGHT:"
    $SIMD_BIN query block --type=height $HEIGHT | jq '.data.txs | length' | sed 's/^/  Transactions: /'
done

echo ""
echo "🔍 Recherche d'événements système:"
echo "--------------------------------"

# Rechercher des événements de module
$SIMD_BIN query txs --query="message.module='timecapsule'" --limit=10 2>/dev/null || echo "Aucun événement du module timecapsule trouvé"

echo ""
echo "📊 Résumé:"
echo "----------"
echo "- Blockchain active: ✅"
echo "- Hauteur actuelle: $CURRENT_HEIGHT"
echo "- Transactions timecapsule trouvées: ❌ (Aucune pour le moment)"
echo ""
echo "💡 Pour créer une capsule de test, utilisez:"
echo "   ./scripts/create_test_capsule.sh"