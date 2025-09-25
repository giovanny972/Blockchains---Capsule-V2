package keeper

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/timecapsule/types"
)

// OptimizationEngine handles automatic performance optimization
type OptimizationEngine struct {
	keeper *Keeper
}

// NewOptimizationEngine creates a new optimization engine
func NewOptimizationEngine(keeper *Keeper) *OptimizationEngine {
	return &OptimizationEngine{
		keeper: keeper,
	}
}

// PerformanceMetrics represents detailed performance metrics
type PerformanceMetrics struct {
	AvgCreateTime       float64 `json:"avg_create_time_ms"`
	AvgOpenTime         float64 `json:"avg_open_time_ms"`
	StorageEfficiency   float64 `json:"storage_efficiency_percent"`
	CacheHitRate        float64 `json:"cache_hit_rate_percent"`
	IPFSLatency         float64 `json:"ipfs_latency_ms"`
	BlockchainLatency   float64 `json:"blockchain_latency_ms"`
	ErrorRate           float64 `json:"error_rate_percent"`
	ThroughputTPS       float64 `json:"throughput_tps"`
	OptimizationScore   float64 `json:"optimization_score"`
	Recommendations     []string `json:"recommendations"`
}

// OptimizationRecommendation represents an automatic optimization suggestion
type OptimizationRecommendation struct {
	Type        string                 `json:"type"`
	Priority    string                 `json:"priority"` // low, medium, high, critical
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Impact      string                 `json:"impact"`
	Effort      string                 `json:"effort"` // low, medium, high
	Actions     []string               `json:"actions"`
	Metrics     map[string]interface{} `json:"metrics"`
	EstimatedGain float64              `json:"estimated_gain_percent"`
}

// AutoOptimizer handles automatic system optimization
type AutoOptimizer struct {
	engine *OptimizationEngine
	rules  []OptimizationRule
}

// OptimizationRule defines automatic optimization rules
type OptimizationRule struct {
	Name        string                       `json:"name"`
	Condition   func(*PerformanceMetrics) bool `json:"-"`
	Action      func(context.Context) error    `json:"-"`
	Priority    string                       `json:"priority"`
	Description string                       `json:"description"`
	Enabled     bool                         `json:"enabled"`
}

// GetPerformanceMetrics calculates comprehensive performance metrics
func (oe *OptimizationEngine) GetPerformanceMetrics(ctx context.Context) (*PerformanceMetrics, error) {
	metrics := &PerformanceMetrics{
		Recommendations: []string{},
	}
	
	// Sample data for demonstration (in production, would collect real metrics)
	metrics.AvgCreateTime = 250.0        // ms
	metrics.AvgOpenTime = 180.0          // ms
	metrics.StorageEfficiency = 85.0     // %
	metrics.CacheHitRate = 75.0          // %
	metrics.IPFSLatency = 120.0          // ms
	metrics.BlockchainLatency = 50.0     // ms
	metrics.ErrorRate = 0.1              // %
	metrics.ThroughputTPS = 45.0         // transactions per second
	
	// Calculate optimization score (0-100)
	score := oe.calculateOptimizationScore(metrics)
	metrics.OptimizationScore = score
	
	// Generate recommendations based on metrics
	recommendations := oe.generateRecommendations(metrics)
	for _, rec := range recommendations {
		metrics.Recommendations = append(metrics.Recommendations, rec.Title)
	}
	
	return metrics, nil
}

// calculateOptimizationScore computes an overall optimization score
func (oe *OptimizationEngine) calculateOptimizationScore(metrics *PerformanceMetrics) float64 {
	// Weighted scoring algorithm
	weights := map[string]float64{
		"storage_efficiency": 0.20,
		"cache_hit_rate":    0.15,
		"latency":           0.25,
		"error_rate":        0.15,
		"throughput":        0.25,
	}
	
	// Normalize metrics to 0-100 scale
	storageScore := metrics.StorageEfficiency
	cacheScore := metrics.CacheHitRate
	
	// Latency score (lower is better)
	avgLatency := (metrics.IPFSLatency + metrics.BlockchainLatency) / 2
	latencyScore := math.Max(0, 100 - (avgLatency / 10)) // 1000ms = 0 score
	
	// Error rate score (lower is better)
	errorScore := math.Max(0, 100 - (metrics.ErrorRate * 100))
	
	// Throughput score (normalized to typical Cosmos TPS)
	throughputScore := math.Min(100, (metrics.ThroughputTPS / 1000) * 100)
	
	// Calculate weighted average
	totalScore := storageScore*weights["storage_efficiency"] +
		cacheScore*weights["cache_hit_rate"] +
		latencyScore*weights["latency"] +
		errorScore*weights["error_rate"] +
		throughputScore*weights["throughput"]
	
	return math.Round(totalScore*100) / 100
}

// generateRecommendations creates optimization recommendations based on metrics
func (oe *OptimizationEngine) generateRecommendations(metrics *PerformanceMetrics) []*OptimizationRecommendation {
	var recommendations []*OptimizationRecommendation
	
	// Storage efficiency recommendations
	if metrics.StorageEfficiency < 80 {
		recommendations = append(recommendations, &OptimizationRecommendation{
			Type:        "storage",
			Priority:    "high",
			Title:       "Optimize Storage Allocation",
			Description: "Storage efficiency is below 80%. Consider implementing better compression and data deduplication.",
			Impact:      "Reduce storage costs by 15-25%",
			Effort:      "medium",
			Actions: []string{
				"Enable compression for large data",
				"Implement data deduplication",
				"Optimize IPFS pinning strategy",
			},
			EstimatedGain: 20.0,
		})
	}
	
	// Cache hit rate recommendations
	if metrics.CacheHitRate < 70 {
		recommendations = append(recommendations, &OptimizationRecommendation{
			Type:        "caching",
			Priority:    "medium",
			Title:       "Improve Cache Strategy",
			Description: "Cache hit rate is below 70%. Implementing smarter caching can improve performance.",
			Impact:      "Reduce response times by 30-40%",
			Effort:      "low",
			Actions: []string{
				"Implement predictive caching",
				"Increase cache size for hot data",
				"Add cache warming for popular capsules",
			},
			EstimatedGain: 15.0,
		})
	}
	
	// Latency recommendations
	avgLatency := (metrics.IPFSLatency + metrics.BlockchainLatency) / 2
	if avgLatency > 100 {
		recommendations = append(recommendations, &OptimizationRecommendation{
			Type:        "latency",
			Priority:    "high",
			Title:       "Reduce Network Latency",
			Description: fmt.Sprintf("Average latency is %.1fms. Consider optimizing network connections.", avgLatency),
			Impact:      "Improve user experience significantly",
			Effort:      "high",
			Actions: []string{
				"Deploy edge nodes closer to users",
				"Optimize IPFS gateway selection",
				"Implement connection pooling",
			},
			EstimatedGain: 25.0,
		})
	}
	
	// Error rate recommendations
	if metrics.ErrorRate > 0.5 {
		recommendations = append(recommendations, &OptimizationRecommendation{
			Type:        "reliability",
			Priority:    "critical",
			Title:       "Improve System Reliability",
			Description: fmt.Sprintf("Error rate is %.2f%%. This needs immediate attention.", metrics.ErrorRate),
			Impact:      "Prevent data loss and improve user trust",
			Effort:      "high",
			Actions: []string{
				"Implement better error handling",
				"Add retry mechanisms",
				"Improve monitoring and alerting",
			},
			EstimatedGain: 30.0,
		})
	}
	
	// Throughput recommendations
	if metrics.ThroughputTPS < 30 {
		recommendations = append(recommendations, &OptimizationRecommendation{
			Type:        "throughput",
			Priority:    "medium",
			Title:       "Increase Transaction Throughput",
			Description: fmt.Sprintf("Throughput is %.1f TPS. Consider batch processing optimizations.", metrics.ThroughputTPS),
			Impact:      "Handle more users simultaneously",
			Effort:      "medium",
			Actions: []string{
				"Implement batch processing",
				"Optimize transaction validation",
				"Use parallel processing where possible",
			},
			EstimatedGain: 40.0,
		})
	}
	
	// Sort by priority and estimated gain
	sort.Slice(recommendations, func(i, j int) bool {
		priorityOrder := map[string]int{"critical": 4, "high": 3, "medium": 2, "low": 1}
		pi, pj := priorityOrder[recommendations[i].Priority], priorityOrder[recommendations[j].Priority]
		if pi != pj {
			return pi > pj
		}
		return recommendations[i].EstimatedGain > recommendations[j].EstimatedGain
	})
	
	return recommendations
}

// GetOptimizationRecommendations returns detailed optimization recommendations
func (oe *OptimizationEngine) GetOptimizationRecommendations(ctx context.Context) ([]*OptimizationRecommendation, error) {
	metrics, err := oe.GetPerformanceMetrics(ctx)
	if err != nil {
		return nil, err
	}
	
	return oe.generateRecommendations(metrics), nil
}

// AutoOptimize performs automatic optimizations based on current metrics
func (oe *OptimizationEngine) AutoOptimize(ctx context.Context) error {
	metrics, err := oe.GetPerformanceMetrics(ctx)
	if err != nil {
		return err
	}
	
	optimizations := []struct {
		name      string
		condition func() bool
		action    func() error
	}{
		{
			name: "compress_large_data",
			condition: func() bool {
				return metrics.StorageEfficiency < 80
			},
			action: func() error {
				return oe.optimizeDataCompression(ctx)
			},
		},
		{
			name: "cleanup_expired_cache",
			condition: func() bool {
				return metrics.CacheHitRate < 70
			},
			action: func() error {
				return oe.cleanupExpiredCache(ctx)
			},
		},
		{
			name: "optimize_ipfs_connections",
			condition: func() bool {
				return metrics.IPFSLatency > 150
			},
			action: func() error {
				return oe.optimizeIPFSConnections(ctx)
			},
		},
	}
	
	optimizationsApplied := 0
	for _, opt := range optimizations {
		if opt.condition() {
			if err := opt.action(); err != nil {
				// Log error but continue with other optimizations
				fmt.Printf("Auto-optimization '%s' failed: %v\n", opt.name, err)
			} else {
				optimizationsApplied++
			}
		}
	}
	
	// Emit optimization event
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"auto_optimization_completed",
			sdk.NewAttribute("optimizations_applied", fmt.Sprintf("%d", optimizationsApplied)),
			sdk.NewAttribute("optimization_score", fmt.Sprintf("%.1f", metrics.OptimizationScore)),
		),
	)
	
	return nil
}

// optimizeDataCompression performs automatic data compression optimization
func (oe *OptimizationEngine) optimizeDataCompression(ctx context.Context) error {
	// In production, this would:
	// - Analyze data patterns
	// - Apply better compression algorithms
	// - Deduplicate similar data
	
	fmt.Println("Applying data compression optimizations...")
	return nil
}

// cleanupExpiredCache removes expired cache entries to improve hit rates
func (oe *OptimizationEngine) cleanupExpiredCache(ctx context.Context) error {
	// In production, this would:
	// - Remove expired cache entries
	// - Reorganize cache for better performance
	// - Pre-warm cache with popular data
	
	fmt.Println("Cleaning up expired cache entries...")
	return nil
}

// optimizeIPFSConnections improves IPFS connection performance
func (oe *OptimizationEngine) optimizeIPFSConnections(ctx context.Context) error {
	// In production, this would:
	// - Select optimal IPFS gateways
	// - Optimize connection pooling
	// - Pin frequently accessed content
	
	fmt.Println("Optimizing IPFS connections...")
	return nil
}

// GetStorageOptimizationReport provides detailed storage analysis
func (oe *OptimizationEngine) GetStorageOptimizationReport(ctx context.Context) (map[string]interface{}, error) {
	report := make(map[string]interface{})
	
	var totalSize, blockchainSize, ipfsSize int64
	var blockchainCount, ipfsCount int64
	
	err := oe.keeper.capsules.Walk(ctx, nil, func(id uint64, capsule types.TimeCapsule) (bool, error) {
		totalSize += capsule.DataSize
		
		if capsule.StorageType == "blockchain" {
			blockchainSize += capsule.DataSize
			blockchainCount++
		} else {
			ipfsSize += capsule.DataSize
			ipfsCount++
		}
		
		return false, nil
	})
	
	if err != nil {
		return nil, err
	}
	
	report["total_data_size"] = totalSize
	report["blockchain_size"] = blockchainSize
	report["ipfs_size"] = ipfsSize
	report["blockchain_count"] = blockchainCount
	report["ipfs_count"] = ipfsCount
	
	if totalSize > 0 {
		report["blockchain_percentage"] = float64(blockchainSize) / float64(totalSize) * 100
		report["ipfs_percentage"] = float64(ipfsSize) / float64(totalSize) * 100
	}
	
	// Calculate cost efficiency
	report["storage_efficiency_score"] = oe.calculateStorageEfficiency(blockchainSize, ipfsSize)
	
	// Recommendations
	recommendations := []string{}
	if blockchainSize > ipfsSize*2 {
		recommendations = append(recommendations, "Consider moving large files to IPFS to reduce blockchain bloat")
	}
	if ipfsCount > blockchainCount*10 {
		recommendations = append(recommendations, "Many small files in IPFS - consider bundling for efficiency")
	}
	
	report["recommendations"] = recommendations
	
	return report, nil
}

// calculateStorageEfficiency computes storage efficiency score
func (oe *OptimizationEngine) calculateStorageEfficiency(blockchainSize, ipfsSize int64) float64 {
	const blockchainCostPer MB = 10.0 // Relative cost units
	const ipfsCostPerMB = 1.0
	
	blockchainMB := float64(blockchainSize) / (1024 * 1024)
	ipfsMB := float64(ipfsSize) / (1024 * 1024)
	
	totalCost := blockchainMB*blockchainCostPer MB + ipfsMB*ipfsCostPerMB
	optimizedCost := math.Min(blockchainMB, 1)*blockchainCostPer MB + (blockchainMB+ipfsMB-math.Min(blockchainMB, 1))*ipfsCostPerMB
	
	if totalCost == 0 {
		return 100.0
	}
	
	efficiency := (1 - (totalCost-optimizedCost)/totalCost) * 100
	return math.Max(0, math.Min(100, efficiency))
}

// SchedulePeriodicOptimization sets up automatic periodic optimization
func (oe *OptimizationEngine) SchedulePeriodicOptimization(ctx context.Context, intervalHours int) error {
	// In production, this would integrate with a job scheduler
	// For now, emit an event that external systems can use
	
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			"periodic_optimization_scheduled",
			sdk.NewAttribute("interval_hours", fmt.Sprintf("%d", intervalHours)),
			sdk.NewAttribute("timestamp", sdkCtx.BlockTime().Format(time.RFC3339)),
		),
	)
	
	return nil
}

// PredictPerformanceImpact predicts the impact of potential changes
func (oe *OptimizationEngine) PredictPerformanceImpact(ctx context.Context, changes map[string]interface{}) (map[string]float64, error) {
	predictions := make(map[string]float64)
	
	// Simple prediction model (in production, would use ML models)
	if cacheIncrease, ok := changes["cache_size_increase"]; ok {
		if increase, ok := cacheIncrease.(float64); ok {
			predictions["response_time_improvement"] = increase * 0.3 // 30% improvement per 100% cache increase
			predictions["storage_cost_increase"] = increase * 0.1     // 10% storage cost increase
		}
	}
	
	if compressionLevel, ok := changes["compression_level"]; ok {
		if level, ok := compressionLevel.(float64); ok {
			predictions["storage_savings"] = level * 0.4     // 40% savings at max compression
			predictions["cpu_overhead"] = level * 0.2        // 20% CPU overhead at max compression
		}
	}
	
	if ipfsOptimization, ok := changes["ipfs_optimization"]; ok {
		if enabled, ok := ipfsOptimization.(bool); ok && enabled {
			predictions["ipfs_latency_reduction"] = 25.0 // 25% latency reduction
			predictions["ipfs_reliability_improvement"] = 15.0 // 15% reliability improvement
		}
	}
	
	return predictions, nil
}