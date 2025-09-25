# Guide d'intégration du module Time Capsule

## 🎯 Objectif
Ce guide explique comment intégrer le module Time Capsule dans votre blockchain Cosmos SDK pour créer un système de capsules temporelles sécurisées.

## 📦 Architecture du Système

### Composants principaux
1. **Module Time Capsule** - Logique principale de gestion des capsules
2. **Système de cryptage AES-256-GCM** - Chiffrement des données
3. **Shamir Secret Sharing** - Distribution des clés entre les masternodes
4. **Smart Contracts de conditions** - Gestion des conditions d'accès
5. **Système de masternodes** - Stockage décentralisé des parts de clés

## 🔧 Intégration dans simapp/app.go

### 1. Ajout des imports
```go
import (
    // Ajoutez ces imports dans simapp/app.go
    timecapsulekeeper "github.com/cosmos/cosmos-sdk/x/timecapsule/keeper"
    timecapsulemodule "github.com/cosmos/cosmos-sdk/x/timecapsule"
    timecapsuletypes "github.com/cosmos/cosmos-sdk/x/timecapsule/types"
)
```

### 2. Ajout du keeper dans la structure SimApp
```go
type SimApp struct {
    // ... autres keepers existants
    
    // Ajoutez le TimeCapsule keeper
    TimeCapsuleKeeper timecapsulekeeper.Keeper
    
    // ... reste de la structure
}
```

### 3. Initialisation du keeper dans NewSimApp
```go
func NewSimApp(...) *SimApp {
    // ... code existant ...
    
    // Ajoutez la clé de store pour timecapsule
    keys := sdk.NewKVStoreKeys(
        // ... autres clés existantes
        timecapsuletypes.StoreKey,
    )
    
    // Initialisez le TimeCapsule keeper
    app.TimeCapsuleKeeper = timecapsulekeeper.NewKeeper(
        appCodec,
        runtime.NewKVStoreService(keys[timecapsuletypes.StoreKey]),
        logger,
        app.BankKeeper,
        app.AccountKeeper,
    )
    
    // Ajoutez au module manager
    app.ModuleManager = module.NewManager(
        // ... modules existants
        timecapsulemodule.NewAppModule(
            appCodec,
            app.TimeCapsuleKeeper,
            app.AccountKeeper,
            app.BankKeeper,
            app.interfaceRegistry,
        ),
    )
    
    // Ordre des blocs
    app.ModuleManager.SetOrderBeginBlockers(
        // ... autres modules
        timecapsuletypes.ModuleName,
    )
    
    app.ModuleManager.SetOrderEndBlockers(
        // ... autres modules  
        timecapsuletypes.ModuleName,
    )
    
    // Ordre de génesis
    genesisModuleOrder := []string{
        // ... autres modules
        timecapsuletypes.ModuleName,
    }
    app.ModuleManager.SetOrderInitGenesis(genesisModuleOrder...)
    app.ModuleManager.SetOrderExportGenesis(genesisModuleOrder...)
    
    // ... reste de l'initialisation
    
    return app
}
```

## 🚀 Compilation et lancement

### 1. Ajout des dépendances Go
```bash
cd simapp
go mod tidy
```

### 2. Compilation
```bash
# Depuis la racine du projet
make build

# Ou directement avec Go
cd simapp
go build -o simd ./simd
```

### 3. Initialisation du nœud
```bash
# Initialiser la chaîne
./simd init mychain --chain-id capsule-testnet

# Créer un utilisateur
./simd keys add alice

# Ajouter des fonds au genesis
./simd genesis add-genesis-account alice 1000000stake

# Créer une transaction de genesis
./simd genesis gentx alice 1000000stake --chain-id capsule-testnet

# Collecter les transactions genesis
./simd genesis collect-gentxs

# Lancer le nœud
./simd start
```

## 📖 Utilisation des capsules temporelles

### 1. Création d'une capsule sécurisée
```bash
# Créer une capsule avec verrouillage temporel
./simd tx timecapsule create-capsule \
  --data-file="./my-secret.json" \
  --capsule-type="time_lock" \
  --unlock-time="2025-12-31T23:59:59Z" \
  --threshold=3 \
  --total-shares=5 \
  --recipient="cosmos1..." \
  --from=alice \
  --chain-id=capsule-testnet
```

### 2. Création d'une capsule Dead Man's Switch
```bash
# Capsule qui se déverrouille après inactivité
./simd tx timecapsule create-capsule \
  --data-file="./inheritance.json" \
  --capsule-type="dead_mans_switch" \
  --inactivity-period=2592000 \
  --threshold=2 \
  --total-shares=3 \
  --recipient="cosmos1heir..." \
  --from=alice
```

### 3. Ouverture d'une capsule
```bash
# Ouvrir une capsule (nécessite les parts de clés)
./simd tx timecapsule open-capsule 1 \
  --key-shares="share1,share2,share3" \
  --from=recipient
```

### 4. Consultation des capsules
```bash
# Lister les capsules d'un utilisateur
./simd query timecapsule list-capsules cosmos1...

# Détails d'une capsule spécifique
./simd query timecapsule capsule 1

# Statistiques globales
./simd query timecapsule stats
```

## 🔒 Sécurité et bonnes pratiques

### 1. Gestion des clés
- Les clés de chiffrement ne sont jamais stockées en clair
- Utilisation de Shamir Secret Sharing avec seuil configurable
- Distribution automatique des parts entre les masternodes

### 2. Types de capsules supportées
- **SAFE**: Stockage sécurisé avec accès propriétaire
- **TIME_LOCK**: Déverrouillage à une date précise
- **CONDITIONAL**: Basé sur des conditions de smart contract
- **MULTI_SIG**: Nécessite plusieurs signatures
- **DEAD_MANS_SWITCH**: Déverrouillage après inactivité

### 3. Paramètres de sécurité
```bash
# Paramètres configurables
./simd query timecapsule params

# Exemple de sortie:
{
  "max_data_size": "1048576",      # 1MB max
  "max_capsule_duration": "8760h", # 1 an max
  "min_threshold": "2",            # Seuil minimum
  "max_shares": "10",              # Parts maximum
  "creation_fee": "100000stake",   # Frais de création
  "maintenance_fee": "10000stake"  # Frais de maintenance
}
```

## 🧪 Tests et validation

### 1. Tests unitaires
```bash
# Lancer les tests du module
go test ./x/timecapsule/...

# Tests de cryptographie
go test ./x/timecapsule/crypto/...

# Tests de keeper
go test ./x/timecapsule/keeper/...
```

### 2. Tests d'intégration
```bash
# Tests complets avec simulation
make test-integration

# Tests E2E
make test-e2e
```

## 🔧 Configuration avancée

### 1. Configuration des masternodes
```yaml
# Dans config/app.toml
[timecapsule]
enable_masternode = true
min_stake_amount = "1000000stake"
key_share_encryption = true
backup_frequency = "1h"
```

### 2. Monitoring et métriques
```bash
# Métriques Prometheus disponibles
curl http://localhost:26660/metrics | grep timecapsule

# Exemples:
# timecapsule_active_capsules
# timecapsule_total_data_encrypted
# timecapsule_key_shares_distributed
```

## 📚 API et requêtes

### Endpoints REST
- `GET /cosmos/timecapsule/v1/capsules` - Liste des capsules
- `GET /cosmos/timecapsule/v1/capsule/{id}` - Détails capsule
- `GET /cosmos/timecapsule/v1/user/{address}/capsules` - Capsules utilisateur
- `POST /cosmos/timecapsule/v1/capsule` - Créer capsule

### Endpoints gRPC
- `cosmos.timecapsule.v1.Query/Capsules`
- `cosmos.timecapsule.v1.Query/Capsule`
- `cosmos.timecapsule.v1.Query/UserCapsules`

## 🚨 Considérations de production

### 1. Sauvegarde et récupération
- Sauvegarde régulière des parts de clés
- Mécanisme de récupération en cas de perte de masternodes
- Système de redondance géographique

### 2. Mise à l'échelle
- Support pour des milliers de capsules simultanées
- Optimisation des performances pour de gros volumes de données
- Cache intelligent pour les requêtes fréquentes

### 3. Conformité réglementaire
- Chiffrement conforme aux standards FIPS 140-2
- Audit trail complet des accès
- Support pour la suppression de données (RGPD)

---

✅ **Le système de capsules temporelles est maintenant prêt à être déployé !**

Ce système offre une sécurité de niveau militaire pour protéger les données sensibles sur la blockchain, avec des mécanismes flexibles d'accès conditionnel inspirés de l'architecture Ternoa.