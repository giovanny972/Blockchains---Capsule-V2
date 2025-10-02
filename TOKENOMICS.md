# 💰 CAPSULE NETWORK - TOKENOMICS

## 📊 Vue d'ensemble

**Capsule Network** utilise un modèle à double token optimisé pour la gouvernance décentralisée et l'utilité du réseau.

---

## 🪙 LES DEUX TOKENS

### **1. STAKE - Token de Gouvernance & Sécurité**
- **Symbole**: `stake`
- **Supply Maximum**: 1 000 000 000 (1 milliard)
- **Usage**:
  - Proof-of-Stake (sécurisation du réseau)
  - Gouvernance on-chain
  - Récompenses de validation

### **2. UCAPS - Token Utilitaire**
- **Symbole**: `ucaps`
- **Supply Maximum**: 100 000 000 (100 millions)
- **Usage**:
  - Frais de transaction (gas)
  - Création de capsules temporelles
  - Paiement des services premium

---

## 📈 RÉPARTITION INITIALE AU GENESIS

### **STAKE (1 000 000 000 tokens)**

| Allocation | Montant | % | Vesting | Usage |
|------------|---------|---|---------|-------|
| **Validator Genesis** | 500 000 000 | 50% | Aucun | Sécurisation réseau + réserve |
| **Équipe & Fondateurs** | 150 000 000 | 15% | 4 ans | Développement à long terme |
| **Trésorerie DAO** | 200 000 000 | 20% | Gouvernance | Financement communautaire |
| **Partenaires Stratégiques** | 50 000 000 | 5% | 2 ans | Écosystème & intégrations |
| **Récompenses Staking** | 100 000 000 | 10% | Émission progressive | APY validateurs (20 ans) |

### **UCAPS (100 000 000 tokens)**

| Allocation | Montant | % | Vesting | Usage |
|------------|---------|---|---------|-------|
| **Validator Genesis** | 10 000 000 | 10% | Aucun | Bootstrap réseau |
| **Programme Beta** | 20 000 000 | 20% | 6 mois | Early adopters & testeurs |
| **Réserve Liquidité** | 30 000 000 | 30% | Lock 1 an | DEX & market making |
| **Trésorerie Opérationnelle** | 25 000 000 | 25% | Aucun | Opérations & marketing |
| **Airdrops & Incentives** | 15 000 000 | 15% | Campagnes | Acquisition utilisateurs |

---

## ⚙️ PARAMÈTRES DE STAKING

### **Validator Genesis (Configuration Actuelle)**

```bash
# Supply allouée au validateur genesis
Total Allocation:     500 000 000 STAKE
Staké Initialement:   200 000 000 STAKE (40%)
Libre (Réserve):      300 000 000 STAKE (60%)
```

### **Métriques de Sécurité**

| Métrique | Valeur | Objectif |
|----------|--------|----------|
| **Stake Initial** | 200M STAKE | Sécuriser le réseau au lancement |
| **Voting Power** | 100% (nœud unique) | Décentralisation progressive |
| **Min Self-Delegation** | 1M STAKE | Protection contre dilution |
| **Unbonding Period** | 21 jours | Standard Cosmos SDK |

### **Récompenses de Validation**

```
APY Cible:              12-18% annuel
Inflation Initiale:     7% par an
Inflation Finale:       2% par an (après 10 ans)
Distribution:           Par bloc (~5 secondes)
Commission Validator:   5% des récompenses
```

---

## 💸 FRAIS DE TRANSACTION (GAS)

### **Modèle de Frais**

| Opération | Coût (UCAPS) | Équivalent USD* |
|-----------|--------------|-----------------|
| **Transfer simple** | 0.001 | ~$0.001 |
| **Création capsule** | 0.1 - 10 | ~$0.10 - $10 |
| **Délégation stake** | 0.01 | ~$0.01 |
| **Vote gouvernance** | 0.05 | ~$0.05 |
| **Création validateur** | 100 | ~$100 |

*_Estimations basées sur 1 UCAPS = $1 (à ajuster selon marché)_

### **Gas Economics**

```
Min Gas Price:     0.001 ucaps
Default Gas Limit: 200 000 units
Max Gas Limit:     10 000 000 units
Fee Burn Rate:     20% (déflationniste)
```

---

## 🔄 MÉCANISME D'INFLATION

### **STAKE - Inflation Dynamique**

```python
# Formule d'inflation annuelle
if bonded_ratio < 67%:
    inflation = 7% + (67% - bonded_ratio) * 0.1
elif bonded_ratio > 67%:
    inflation = max(7% - (bonded_ratio - 67%) * 0.1, 2%)
else:
    inflation = 7%

# Réduction progressive
inflation_reduction = 0.5% par an (jusqu'à 2% minimum)
```

**Objectif**: Maintenir ~67% des tokens stakés pour sécurité optimale

### **UCAPS - Supply Fixe avec Burn**

```
Supply Maximum:    100 000 000 (fixe)
Pas d'inflation:   ✅
Burn des frais:    20% par transaction
Supply Finale:     Déflationniste (tend vers 80M à long terme)
```

---

## 📊 MÉTRIQUES CLÉS (ÉTAT ACTUEL)

### **Snapshot Genesis**

```yaml
Chain ID:          capsule-mainnet-1
Block Height:      0 (genesis)
Total Validators:  1

STAKE Token:
  Total Supply:    1 000 000 000
  Circulating:     500 000 000 (50%)
  Bonded:          200 000 000 (20% du total)
  Bonded Ratio:    40% (du circulant)

UCAPS Token:
  Total Supply:    100 000 000
  Circulating:     10 000 000 (10%)
  Burned:          0
```

### **Objectifs à 6 Mois**

```yaml
Target Validators:    25-50 actifs
Target Bonded Ratio:  67% (670M STAKE)
Target Circulation:
  - STAKE: 850M (85%)
  - UCAPS: 60M (60%)
```

---

## 🎯 STRATÉGIE DE DÉCENTRALISATION

### **Phase 1: Bootstrap (Mois 1-3)**
- ✅ Validateur genesis opérationnel
- 🎯 Recrutement 10 validateurs partenaires
- 🎯 Distribution initiale STAKE (airdrops, incentives)
- 🎯 Lancement programme beta UCAPS

### **Phase 2: Expansion (Mois 4-12)**
- 🎯 50+ validateurs actifs
- 🎯 Bonded ratio >67%
- 🎯 Lancement DEX pour liquidité UCAPS
- 🎯 Activation gouvernance on-chain

### **Phase 3: Maturité (Année 2+)**
- 🎯 100+ validateurs géographiquement distribués
- 🎯 Décentralisation complète (aucun validateur >10% VP)
- 🎯 Trésorerie DAO autonome
- 🎯 Émissions STAKE réduites à 2-3% annuel

---

## 🔐 SÉCURITÉ & GOUVERNANCE

### **Paramètres de Gouvernance**

```yaml
Quorum Minimum:        33.4% des tokens stakés
Seuil d'Approbation:   50% des votes (hors abstention)
Seuil de Veto:         33.4% pour bloquer
Période de Vote:       7 jours
Période de Dépôt:      3 jours
Dépôt Minimum:         1000 STAKE
```

### **Protection Anti-Centralisation**

- **Max Commission Rate**: 20%
- **Max Commission Change**: 1% par jour
- **Slashing Downtime**: 0.01% après 24h offline
- **Slashing Double-Sign**: 5% + jail permanent
- **Min Self-Bond**: 1M STAKE pour validateurs

---

## 💡 COMPARAISON AVEC D'AUTRES CHAINS

| Chain | Supply Token | Inflation | Bonded Target | Unbonding |
|-------|--------------|-----------|---------------|-----------|
| **Capsule Network** | 1B STAKE | 2-7% | 67% | 21 jours |
| **Cosmos Hub (ATOM)** | 390M | 7-20% | 67% | 21 jours |
| **Osmosis (OSMO)** | 1B | Décroissant | 67% | 14 jours |
| **Juno (JUNO)** | 185M | 40% → 10% | 67% | 28 jours |

**Avantages Capsule**:
- ✅ Supply modérée (ni trop rare, ni trop dilué)
- ✅ Inflation contrôlée et prévisible
- ✅ Double token (séparation gouvernance/utilité)
- ✅ Mécanisme déflationniste sur UCAPS

---

## 📞 CONTACT & RESSOURCES

- **Documentation**: https://docs.capsulenetwork.io
- **Explorer**: https://explorer.capsulenetwork.io
- **GitHub**: https://github.com/capsule-network
- **Discord**: https://discord.gg/capsulenetwork

---

**Dernière mise à jour**: 2025-10-02
**Version**: 1.0.0
**Chain**: capsule-mainnet-1
