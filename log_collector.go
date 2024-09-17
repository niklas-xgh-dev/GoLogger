package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

type SecurityLog struct {
	Timestamp time.Time       `json:"timestamp"`
	LogData   json.RawMessage `json:"log_data"`
}

type LogStorage interface {
	CreateTableIfNotExists() error
	InsertLog(log SecurityLog) error
}

type PostgresStorage struct {
	db *sql.DB
}

func (p *PostgresStorage) CreateTableIfNotExists() error {
	_, err := p.db.Exec(`
		CREATE TABLE IF NOT EXISTS security_logs (
			id SERIAL PRIMARY KEY,
			timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
			log_data JSONB NOT NULL
		)
	`)
	return err
}

func (p *PostgresStorage) InsertLog(log SecurityLog) error {
	_, err := p.db.Exec(`
		INSERT INTO security_logs (timestamp, log_data)
		VALUES ($1, $2)
	`, log.Timestamp, log.LogData)
	return err
}

func collectSecurityLogs() (map[string]interface{}, error) {
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil {
		return nil, fmt.Errorf("error getting CPU percent: %v", err)
	}
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("error getting memory info: %v", err)
	}
	diskInfo, err := disk.Usage("/")
	if err != nil {
		return nil, fmt.Errorf("error getting disk usage: %v", err)
	}

	openPorts, _ := getOpenPorts()
	runningServices, _ := getRunningServices()
	systemUsers, _ := getSystemUsers()
	recentLogins, _ := getRecentLogins()
	failedLogins, _ := getFailedLogins()
	networkConnections, _ := net.Connections("all")
	suspiciousProcesses, _ := getSuspiciousProcesses()

	return map[string]interface{}{
		"timestamp":            time.Now(),
		"cpu_percent":          cpuPercent[0],
		"memory_percent":       memInfo.UsedPercent,
		"disk_percent":         diskInfo.UsedPercent,
		"open_ports":           openPorts,
		"running_services":     runningServices,
		"system_users":         systemUsers,
		"recent_logins":        recentLogins,
		"failed_logins":        failedLogins,
		"network_connections":  networkConnections,
		"suspicious_processes": suspiciousProcesses,
	}, nil
}

func getOpenPorts() ([]int, error) {
	cmd := exec.Command("netstat", "-tuln")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var ports []int
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "LISTEN") {
			fields := strings.Fields(line)
			if len(fields) > 3 {
				addrPort := strings.Split(fields[3], ":")
				if port, err := strconv.Atoi(addrPort[len(addrPort)-1]); err == nil {
					ports = append(ports, port)
				}
			}
		}
	}
	return ports, nil
}

func getRunningServices() ([]string, error) {
	if runtime.GOOS == "darwin" {
		cmd := exec.Command("launchctl", "list")
		output, err := cmd.Output()
		if err != nil {
			return nil, err
		}
		return strings.Split(string(output), "\n"), nil
	} else if runtime.GOOS == "linux" {
		cmd := exec.Command("systemctl", "list-units", "--type=service")
		output, err := cmd.Output()
		if err != nil {
			return nil, err
		}
		return strings.Split(string(output), "\n"), nil
	}
	return nil, fmt.Errorf("unsupported operating system")
}

func getSystemUsers() ([]string, error) {
	data, err := os.ReadFile("/etc/passwd")
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(data), "\n")
	var users []string
	for _, line := range lines {
		if parts := strings.Split(line, ":"); len(parts) > 0 {
			users = append(users, parts[0])
		}
	}
	return users, nil
}

func getRecentLogins() ([]string, error) {
	cmd := exec.Command("last", "-n", "10")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return strings.Split(string(output), "\n"), nil
}

func getFailedLogins() (int, error) {
	// Placeholder: Implement based on your OS and logging system
	return 0, nil
}

func getSuspiciousProcesses() ([]string, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var suspicious []string
	for _, p := range processes {
		name, err := p.Name()
		if err != nil {
			continue
		}
		if strings.Contains(strings.ToLower(name), "suspicious") {
			suspicious = append(suspicious, name)
		}
	}
	return suspicious, nil
}

func main() {
	if err := godotenv.Load("config/.env"); err != nil {
		log.Println("Warning: Error loading .env file")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set. Exiting.")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}
	defer db.Close()

	storage := &PostgresStorage{db: db}

	if err := storage.CreateTableIfNotExists(); err != nil {
		log.Fatalf("Error creating table: %v", err)
	}

	for {
		securityLogData, err := collectSecurityLogs()
		if err != nil {
			log.Printf("Error collecting security logs: %v", err)
			continue
		}

		logJSON, err := json.Marshal(securityLogData)
		if err != nil {
			log.Printf("Error marshaling JSON: %v", err)
			continue
		}

		secLog := SecurityLog{
			Timestamp: time.Now(),
			LogData:   logJSON,
		}

		if err := storage.InsertLog(secLog); err != nil {
			log.Printf("Error inserting log: %v", err)
		} else {
			fmt.Printf("Log inserted at %s\n", secLog.Timestamp)
		}

		time.Sleep(5 * time.Minute)
	}
}
