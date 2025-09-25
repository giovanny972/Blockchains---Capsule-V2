package ipfs

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/x/timecapsule/types"
)

// IPFSManager manages IPFS operations for time capsules
type IPFSManager struct {
	client       *IPFSClient
	config       IPFSConfig
	cache        map[string]*IPFSMetadata
	cacheMutex   sync.RWMutex
	healthStatus bool
	lastHealth   time.Time
}

// NewIPFSManager creates a new IPFS manager
func NewIPFSManager(config IPFSConfig) *IPFSManager {
	return &IPFSManager{
		client:       NewIPFSClient(config),
		config:       config,
		cache:        make(map[string]*IPFSMetadata),
		healthStatus: false,
	}
}

// Initialize initializes the IPFS manager and checks connectivity
func (m *IPFSManager) Initialize(ctx context.Context) error {
	// Health check
	if err := m.client.HealthCheck(ctx); err != nil {
		return fmt.Errorf("IPFS initialization failed: %w", err)
	}

	m.healthStatus = true
	m.lastHealth = time.Now()

	return nil
}

// StoreCapsuleData stores capsule data in IPFS with enhanced metadata
func (m *IPFSManager) StoreCapsuleData(ctx context.Context, capsuleID uint64, encryptedData []byte, capsule *types.TimeCapsule) (*IPFSMetadata, error) {
	// Create comprehensive metadata
	metadata := IPFSMetadata{
		Name:      fmt.Sprintf("capsule_%d_data.enc", capsuleID),
		MimeType:  "application/octet-stream",
		Encrypted: true,
		CapsuleID: capsuleID,
		Pinned:    true,
		Redundancy: 3, // Target 3 IPFS nodes
	}

	// Add expiry time if capsule has unlock time
	if capsule.UnlockTime != nil {
		// Keep data for 30 days after unlock time
		expiryTime := capsule.UnlockTime.Add(30 * 24 * time.Hour)
		metadata.ExpiryTime = &expiryTime
	}

	// Store in IPFS
	result, err := m.client.StoreData(ctx, encryptedData, metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to store capsule data in IPFS: %w", err)
	}

	// Pin data for persistence
	if err := m.client.PinData(ctx, result.Hash); err != nil {
		// Log warning but don't fail - data is still stored
		fmt.Printf("Warning: failed to pin data %s: %v\n", result.Hash, err)
	}

	// Cache metadata
	m.cacheMutex.Lock()
	m.cache[result.Hash] = result
	m.cacheMutex.Unlock()

	return result, nil
}

// RetrieveCapsuleData retrieves and validates capsule data from IPFS
func (m *IPFSManager) RetrieveCapsuleData(ctx context.Context, metadata *IPFSMetadata) ([]byte, error) {
	// Check cache first
	m.cacheMutex.RLock()
	cachedMeta, exists := m.cache[metadata.Hash]
	m.cacheMutex.RUnlock()

	if exists {
		metadata = cachedMeta
	}

	// Retrieve from IPFS
	data, err := m.client.RetrieveData(ctx, *metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve capsule data from IPFS: %w", err)
	}

	return data, nil
}

// CleanupExpiredData removes expired capsule data from IPFS
func (m *IPFSManager) CleanupExpiredData(ctx context.Context) error {
	// Get all pinned content
	pins, err := m.client.ListPins(ctx)
	if err != nil {
		return fmt.Errorf("failed to list pins: %w", err)
	}

	now := time.Now()
	cleanedCount := 0

	for _, hash := range pins {
		// Check if we have metadata for this hash
		m.cacheMutex.RLock()
		metadata, exists := m.cache[hash]
		m.cacheMutex.RUnlock()

		if exists && metadata.ExpiryTime != nil && now.After(*metadata.ExpiryTime) {
			// Unpin expired data
			if err := m.client.UnpinData(ctx, hash); err != nil {
				fmt.Printf("Warning: failed to unpin expired data %s: %v\n", hash, err)
				continue
			}

			// Remove from cache
			m.cacheMutex.Lock()
			delete(m.cache, hash)
			m.cacheMutex.Unlock()

			cleanedCount++
		}
	}

	fmt.Printf("Cleaned up %d expired IPFS entries\n", cleanedCount)
	return nil
}

// GetStorageStats returns IPFS storage statistics
func (m *IPFSManager) GetStorageStats(ctx context.Context) (*StorageStats, error) {
	// Get IPFS repo stats
	stats, err := m.client.GetStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get IPFS stats: %w", err)
	}

	// Count cached metadata
	m.cacheMutex.RLock()
	cachedItems := len(m.cache)
	var totalSize int64
	for _, metadata := range m.cache {
		totalSize += metadata.Size
	}
	m.cacheMutex.RUnlock()

	return &StorageStats{
		TotalSize:       totalSize,
		CachedItems:     cachedItems,
		IPFSRepoSize:    getInt64FromStats(stats, "RepoSize"),
		IPFSStorageMax:  getInt64FromStats(stats, "StorageMax"),
		IPFSNumObjects:  getInt64FromStats(stats, "NumObjects"),
		HealthStatus:    m.healthStatus,
		LastHealthCheck: m.lastHealth,
	}, nil
}

// StorageStats represents IPFS storage statistics
type StorageStats struct {
	TotalSize       int64     `json:"total_size"`
	CachedItems     int       `json:"cached_items"`
	IPFSRepoSize    int64     `json:"ipfs_repo_size"`
	IPFSStorageMax  int64     `json:"ipfs_storage_max"`
	IPFSNumObjects  int64     `json:"ipfs_num_objects"`
	HealthStatus    bool      `json:"health_status"`
	LastHealthCheck time.Time `json:"last_health_check"`
}

// BackupCapsuleData creates redundant copies of capsule data
func (m *IPFSManager) BackupCapsuleData(ctx context.Context, metadata *IPFSMetadata, redundancy int) error {
	// For now, we rely on IPFS network redundancy
	// In the future, this could implement explicit multi-node pinning
	
	// Pin on additional nodes (if available)
	for i := 0; i < redundancy; i++ {
		if err := m.client.PinData(ctx, metadata.Hash); err != nil {
			fmt.Printf("Warning: backup pin attempt %d failed: %v\n", i+1, err)
		}
	}

	return nil
}

// HealthCheck performs periodic health checks
func (m *IPFSManager) HealthCheck(ctx context.Context) error {
	err := m.client.HealthCheck(ctx)
	m.healthStatus = (err == nil)
	m.lastHealth = time.Now()
	
	if err != nil {
		return fmt.Errorf("IPFS health check failed: %w", err)
	}

	return nil
}

// ValidateIPFSHash validates an IPFS hash format
func ValidateIPFSHash(hash string) bool {
	// Basic validation - IPFS hashes typically start with 'Qm' or 'ba'
	if len(hash) < 46 {
		return false
	}
	
	return hash[:2] == "Qm" || hash[:2] == "ba" || hash[:2] == "zb"
}

// Helper function to extract int64 from stats map
func getInt64FromStats(stats map[string]interface{}, key string) int64 {
	if val, exists := stats[key]; exists {
		switch v := val.(type) {
		case float64:
			return int64(v)
		case int64:
			return v
		case int:
			return int64(v)
		}
	}
	return 0
}

// GetMetadataFromCache retrieves metadata from cache
func (m *IPFSManager) GetMetadataFromCache(hash string) (*IPFSMetadata, bool) {
	m.cacheMutex.RLock()
	defer m.cacheMutex.RUnlock()
	
	metadata, exists := m.cache[hash]
	return metadata, exists
}

// UpdateCacheMetadata updates metadata in cache
func (m *IPFSManager) UpdateCacheMetadata(hash string, metadata *IPFSMetadata) {
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()
	
	m.cache[hash] = metadata
}