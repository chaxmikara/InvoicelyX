package handlers

import (
	"context"
	"runtime"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// SystemStats represents system monitoring data
type SystemStats struct {
	DatabaseStats DatabaseStats `json:"database_stats"`
	SystemInfo    SystemInfo    `json:"system_info"`
	APIStats      APIStats      `json:"api_stats"`
	GeneratedAt   time.Time     `json:"generated_at"`
}

type DatabaseStats struct {
	TotalUsers    int64 `json:"total_users"`
	ActiveUsers   int64 `json:"active_users"`
	InactiveUsers int64 `json:"inactive_users"`
	TotalInvoices int64 `json:"total_invoices"`
}

type SystemInfo struct {
	GoVersion     string `json:"go_version"`
	NumGoroutines int    `json:"num_goroutines"`
	MemAllocMB    uint64 `json:"memory_alloc_mb"`
	MemSysMB      uint64 `json:"memory_sys_mb"`
	NumGC         uint32 `json:"num_gc"`
}

type APIStats struct {
	Uptime         string    `json:"uptime"`
	ServerStarted  time.Time `json:"server_started"`
	TotalEndpoints int       `json:"total_endpoints"`
}

var serverStartTime = time.Now()

// GetSystemStats - Admin endpoint for system monitoring
func (h *UserHandler) GetSystemStats(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Database statistics
	usersCollection := h.DB.Collection("users")
	invoicesCollection := h.DB.Collection("invoices")

	totalUsers, _ := usersCollection.CountDocuments(ctx, bson.M{})
	activeUsers, _ := usersCollection.CountDocuments(ctx, bson.M{"isActive": true})
	inactiveUsers, _ := usersCollection.CountDocuments(ctx, bson.M{"isActive": false})
	totalInvoices, _ := invoicesCollection.CountDocuments(ctx, bson.M{})

	// System information
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	systemStats := SystemStats{
		DatabaseStats: DatabaseStats{
			TotalUsers:    totalUsers,
			ActiveUsers:   activeUsers,
			InactiveUsers: inactiveUsers,
			TotalInvoices: totalInvoices,
		},
		SystemInfo: SystemInfo{
			GoVersion:     runtime.Version(),
			NumGoroutines: runtime.NumGoroutine(),
			MemAllocMB:    memStats.Alloc / 1024 / 1024,
			MemSysMB:      memStats.Sys / 1024 / 1024,
			NumGC:         memStats.NumGC,
		},
		APIStats: APIStats{
			Uptime:         time.Since(serverStartTime).String(),
			ServerStarted:  serverStartTime,
			TotalEndpoints: 15, //....
		},
		GeneratedAt: time.Now(),
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "System statistics retrieved successfully",
		"data":    systemStats,
	})
}

// GetSystemHealth - Admin endpoint for health check
func (h *UserHandler) GetSystemHealth(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test database connectivity
	dbHealthy := true
	var dbError string
	err := h.DB.Client().Ping(ctx, nil)
	if err != nil {
		dbHealthy = false
		dbError = err.Error()
	}

	// Memory check
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	memoryUsageMB := memStats.Alloc / 1024 / 1024

	// Determine overall health
	overallHealth := "healthy"
	if !dbHealthy || memoryUsageMB > 500 { // 500MB threshold
		overallHealth = "unhealthy"
	}

	healthData := fiber.Map{
		"status":    overallHealth,
		"timestamp": time.Now(),
		"uptime":    time.Since(serverStartTime).String(),
		"checks": fiber.Map{
			"database": fiber.Map{
				"status": map[bool]string{true: "healthy", false: "unhealthy"}[dbHealthy],
				"error":  dbError,
			},
			"memory": fiber.Map{
				"status":    map[bool]string{true: "healthy", false: "warning"}[memoryUsageMB < 500],
				"usage_mb":  memoryUsageMB,
				"threshold": 500,
			},
			"goroutines": fiber.Map{
				"count":  runtime.NumGoroutine(),
				"status": "healthy",
			},
		},
	}

	statusCode := fiber.StatusOK
	if overallHealth == "unhealthy" {
		statusCode = fiber.StatusServiceUnavailable
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"success": overallHealth == "healthy",
		"message": "System health check completed",
		"data":    healthData,
	})
}
