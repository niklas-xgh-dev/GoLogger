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
	Timestamp           time.Time            `json:"timestamp"`
	CPUPercent          float64              `json:"cpu_percent"`
	MemoryPercent       float64              `json:"memory_percent"`
	DiskPercent         float64              `json:"disk_percent"`
	OpenPorts           []int                `json:"open_ports"`
	RunningServices     []string             `json:"running_services"`
	SystemUsers         []string             `json:"system_users"`
	RecentLogins        []string             `json:"recent_logins"`
	FailedLogins        int                  `json:"failed_logins"`
	NetworkConnections  []net.ConnectionStat `json:"network_connections"`
	SuspiciousProcesses []string             `json:"suspicious_processes"`
}

type LogStorage interface {
	InsertLog(log SecurityLog) error
}

type PostgresStorage struct {
	db *sql.DB
}

func (p *PostgresStorage) InsertLog(log SecurityLog) error {
	jsonData, err := json.Marshal(log)
	if err != nil {
		return err
	}

	_, err = p.db.Exec(`
		INSERT INTO security_logs (timestamp, log_data)
		VALUES ($1, $2)
	`, log.Timestamp, jsonData)
	return err
}

type MockStorage struct {
	logs []SecurityLog
}

func (m *MockStorage) InsertLog(log SecurityLog) error {
	m.logs = append(m.logs, log)
	return nil
}

func collectSecurityLogs() (SecurityLog, error) {
	cpuPercent, _ := cpu.Percent(0, false)
	memInfo, _ := mem.VirtualMemory()
	diskInfo, _ := disk.Usage("/")

	openPorts, _ := getOpenPorts()
	runningServices, _ := getRunningServices()
	systemUsers, _ := getSystemUsers()
	recentLogins, _ := getRecentLogins()
	failedLogins, _ := getFailedLogins()
	networkConnections, _ := net.Connections("all")
	suspiciousProcesses, _ := getSuspiciousProcesses()

	return SecurityLog{
		Timestamp:           time.Now(),
		CPUPercent:          cpuPercent[0],
		MemoryPercent:       memInfo.UsedPercent,
		DiskPercent:         diskInfo.UsedPercent,
		OpenPorts:           openPorts,
		RunningServices:     runningServices,
		SystemUsers:         systemUsers,
		RecentLogins:        recentLogins,
		FailedLogins:        failedLogins,
		NetworkConnections:  networkConnections,
		SuspiciousProcesses: suspiciousProcesses,
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
		// Parse the output to extract service names
		return strings.Split(string(output), "\n"), nil
	} else if runtime.GOOS == "linux" {
		cmd := exec.Command("systemctl", "list-units", "--type=service")
		output, err := cmd.Output()
		if err != nil {
			return nil, err
		}
		// Parse the output to extract service names
		return strings.Split(string(output), "\n"), nil
	}
	return nil, fmt.Errorf("unsupported operating system")
}

func getSystemUsers() ([]string, error) {
	// Read from /etc/passwd on Unix-like systems
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
	// This is a placeholder. The actual implementation would depend on the OS and logging system.
	// For example, on Linux you might parse /var/log/auth.log
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
		// This is a very simplistic check and should be expanded based on your specific criteria
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
		securityLog, err := collectSecurityLogs()
		if err != nil {
			log.Printf("Error collecting security logs: %v", err)
			continue
		}

		logJSON, _ := json.Marshal(securityLog)
		fmt.Printf("Collected security log: %s\n", logJSON)

		if err := storage.InsertLog(securityLog); err != nil {
			log.Printf("Error inserting log: %v", err)
		}

		time.Sleep(5 * time.Minute)
	}
}
