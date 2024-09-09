package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

type SystemLog struct {
	Timestamp     time.Time `json:"timestamp"`
	CPUPercent    float64   `json:"cpu_percent"`
	MemoryPercent float64   `json:"memory_percent"`
	DiskPercent   float64   `json:"disk_percent"`
}

func collectSystemLogs() SystemLog {
	cpuPercent, _ := cpu.Percent(0, false)
	memInfo, _ := mem.VirtualMemory()
	diskInfo, _ := disk.Usage("/")

	return SystemLog{
		Timestamp:     time.Now(),
		CPUPercent:    cpuPercent[0],
		MemoryPercent: memInfo.UsedPercent,
		DiskPercent:   diskInfo.UsedPercent,
	}
}

func insertLog(db *sql.DB, log SystemLog) error {
	_, err := db.Exec(`
		INSERT INTO system_logs (timestamp, cpu_percent, memory_percent, disk_percent)
		VALUES ($1, $2, $3, $4)
	`, log.Timestamp, log.CPUPercent, log.MemoryPercent, log.DiskPercent)
	return err
}

func main() {
	if err := godotenv.Load("config/.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	for {
		systemLog := collectSystemLogs()

		logJSON, _ := json.Marshal(systemLog)
		fmt.Printf("Collected log: %s\n", logJSON)

		if err := insertLog(db, systemLog); err != nil {
			log.Printf("Error inserting log: %v", err)
		}

		time.Sleep(1 * time.Minute)
	}
}
