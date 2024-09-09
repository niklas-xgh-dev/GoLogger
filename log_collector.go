package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemLog struct {
	Timestamp     time.Time `json:"timestamp"`
	CPUPercent    float64   `json:"cpu_percent"`
	MemoryPercent float64   `json:"memory_percent"`
	DiskPercent   float64   `json:"disk_percent"`
}

type LogStorage interface {
	InsertLog(log SystemLog) error
}

type PostgresStorage struct {
	db *sql.DB
}

func (p *PostgresStorage) InsertLog(log SystemLog) error {
	_, err := p.db.Exec(`
		INSERT INTO system_logs (timestamp, cpu_percent, memory_percent, disk_percent)
		VALUES ($1, $2, $3, $4)
	`, log.Timestamp, log.CPUPercent, log.MemoryPercent, log.DiskPercent)
	return err
}

type MockStorage struct {
	logs []SystemLog
}

func (m *MockStorage) InsertLog(log SystemLog) error {
	m.logs = append(m.logs, log)
	return nil
}

func collectSystemLogs() (SystemLog, error) {
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		return SystemLog{}, fmt.Errorf("error collecting CPU usage: %v", err)
	}

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return SystemLog{}, fmt.Errorf("error collecting memory usage: %v", err)
	}

	var diskInfo *disk.UsageStat
	if runtime.GOOS == "windows" {
		diskInfo, err = disk.Usage("C:")
	} else {
		diskInfo, err = disk.Usage("/")
	}
	if err != nil {
		return SystemLog{}, fmt.Errorf("error collecting disk usage: %v", err)
	}

	return SystemLog{
		Timestamp:     time.Now(),
		CPUPercent:    cpuPercent[0],
		MemoryPercent: memInfo.UsedPercent,
		DiskPercent:   diskInfo.UsedPercent,
	}, nil
}

func main() {
	if err := godotenv.Load("config/.env"); err != nil {
		log.Println("Warning: Error loading .env file")
	}

	var storage LogStorage

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("DATABASE_URL not set. Using mock storage.")
		storage = &MockStorage{}
	} else {
		db, err := sql.Open("postgres", dbURL)
		if err != nil {
			log.Fatalf("Error opening database connection: %v", err)
		}
		defer db.Close()
		storage = &PostgresStorage{db: db}
	}

	for {
		systemLog, err := collectSystemLogs()
		if err != nil {
			log.Printf("Error collecting system logs: %v", err)
			continue
		}

		logJSON, _ := json.Marshal(systemLog)
		fmt.Printf("Collected log: %s\n", logJSON)

		if err := storage.InsertLog(systemLog); err != nil {
			log.Printf("Error inserting log: %v", err)
		}

		time.Sleep(1 * time.Minute)
	}
}
