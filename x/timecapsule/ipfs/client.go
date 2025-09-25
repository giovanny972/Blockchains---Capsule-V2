package ipfs

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

// IPFSClient handles interactions with IPFS network
type IPFSClient struct {
	apiURL     string
	gatewayURL string
	client     *http.Client
	timeout    time.Duration
}

// IPFSConfig holds IPFS configuration
type IPFSConfig struct {
	APIEndpoint     string        `json:"api_endpoint"`
	GatewayEndpoint string        `json:"gateway_endpoint"`
	PinningService  string        `json:"pinning_service,omitempty"`
	Timeout         time.Duration `json:"timeout"`
	MaxFileSize     int64         `json:"max_file_size"`
}

// IPFSMetadata contains metadata about stored content
type IPFSMetadata struct {
	Hash         string    `json:"hash"`
	Size         int64     `json:"size"`
	Name         string    `json:"name,omitempty"`
	MimeType     string    `json:"mime_type,omitempty"`
	UploadTime   time.Time `json:"upload_time"`
	DataHash     string    `json:"data_hash"`     // SHA-256 hash of original data
	Encrypted    bool      `json:"encrypted"`
	CapsuleID    uint64    `json:"capsule_id"`
	Pinned       bool      `json:"pinned"`
	Redundancy   int       `json:"redundancy"`    // Number of IPFS nodes storing the data
	ExpiryTime   *time.Time `json:"expiry_time,omitempty"`
}

// IPFSResponse represents response from IPFS API
type IPFSResponse struct {
	Name string `json:"Name"`
	Hash string `json:"Hash"`
	Size string `json:"Size"`
}

// NewIPFSClient creates a new IPFS client
func NewIPFSClient(config IPFSConfig) *IPFSClient {
	return &IPFSClient{
		apiURL:     config.APIEndpoint,
		gatewayURL: config.GatewayEndpoint,
		timeout:    config.Timeout,
		client: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// DefaultIPFSConfig returns default IPFS configuration
func DefaultIPFSConfig() IPFSConfig {
	return IPFSConfig{
		APIEndpoint:     "http://127.0.0.1:5001",
		GatewayEndpoint: "http://127.0.0.1:8080",
		Timeout:         30 * time.Second,
		MaxFileSize:     100 * 1024 * 1024, // 100MB max
	}
}

// StoreData stores encrypted data in IPFS and returns metadata
func (c *IPFSClient) StoreData(ctx context.Context, data []byte, metadata IPFSMetadata) (*IPFSMetadata, error) {
	// Validate data size
	if int64(len(data)) > 100*1024*1024 { // 100MB limit
		return nil, fmt.Errorf("data too large: %d bytes (max 100MB)", len(data))
	}

	// Calculate data hash for integrity
	hash := sha256.Sum256(data)
	metadata.DataHash = hex.EncodeToString(hash[:])
	metadata.Size = int64(len(data))
	metadata.UploadTime = time.Now()

	// Create multipart form for IPFS upload
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	
	// Add file field
	part, err := writer.CreateFormFile("file", metadata.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	
	if _, err := part.Write(data); err != nil {
		return nil, fmt.Errorf("failed to write data: %w", err)
	}
	
	writer.Close()

	// Upload to IPFS
	req, err := http.NewRequestWithContext(ctx, "POST", c.apiURL+"/api/v0/add", &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to upload to IPFS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("IPFS upload failed: %s", string(body))
	}

	// Parse IPFS response
	var ipfsResp IPFSResponse
	if err := json.NewDecoder(resp.Body).Decode(&ipfsResp); err != nil {
		return nil, fmt.Errorf("failed to decode IPFS response: %w", err)
	}

	// Update metadata with IPFS hash
	metadata.Hash = ipfsResp.Hash
	metadata.Pinned = true // Assume pinned by default

	return &metadata, nil
}

// RetrieveData retrieves and validates data from IPFS
func (c *IPFSClient) RetrieveData(ctx context.Context, metadata IPFSMetadata) ([]byte, error) {
	// Retrieve from IPFS gateway
	url := fmt.Sprintf("%s/ipfs/%s", c.gatewayURL, metadata.Hash)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve from IPFS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IPFS retrieval failed: status %d", resp.StatusCode)
	}

	// Read data
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}

	// Validate data integrity
	hash := sha256.Sum256(data)
	expectedHash := metadata.DataHash
	actualHash := hex.EncodeToString(hash[:])
	
	if actualHash != expectedHash {
		return nil, fmt.Errorf("data integrity check failed: expected %s, got %s", expectedHash, actualHash)
	}

	return data, nil
}

// PinData ensures data is pinned on IPFS network for persistence
func (c *IPFSClient) PinData(ctx context.Context, hash string) error {
	url := fmt.Sprintf("%s/api/v0/pin/add?arg=%s", c.apiURL, hash)
	
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create pin request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to pin data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("IPFS pin failed: %s", string(body))
	}

	return nil
}

// UnpinData removes pin from IPFS (allows garbage collection)
func (c *IPFSClient) UnpinData(ctx context.Context, hash string) error {
	url := fmt.Sprintf("%s/api/v0/pin/rm?arg=%s", c.apiURL, hash)
	
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create unpin request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to unpin data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("IPFS unpin failed: %s", string(body))
	}

	return nil
}

// GetStats returns statistics about IPFS node
func (c *IPFSClient) GetStats(ctx context.Context) (map[string]interface{}, error) {
	url := c.apiURL + "/api/v0/stats/repo"
	
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create stats request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IPFS stats failed: status %d", resp.StatusCode)
	}

	var stats map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("failed to decode stats: %w", err)
	}

	return stats, nil
}

// HealthCheck verifies IPFS node connectivity
func (c *IPFSClient) HealthCheck(ctx context.Context) error {
	url := c.apiURL + "/api/v0/version"
	
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create health check request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("IPFS node not reachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("IPFS node unhealthy: status %d", resp.StatusCode)
	}

	return nil
}

// ListPins returns all pinned content
func (c *IPFSClient) ListPins(ctx context.Context) ([]string, error) {
	url := c.apiURL + "/api/v0/pin/ls"
	
	req, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create list pins request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to list pins: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IPFS list pins failed: status %d", resp.StatusCode)
	}

	var result map[string]map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode pins list: %w", err)
	}

	var hashes []string
	for hash := range result["Keys"] {
		hashes = append(hashes, hash)
	}

	return hashes, nil
}