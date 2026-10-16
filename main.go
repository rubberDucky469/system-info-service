package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type StatusResponse struct {
	Hostname      string  `json:"hostname"`
	UptimeSeconds float64 `json:"uptime_seconds"`
	UptimeHuman   string  `json:"uptime_human"`
}

func getUptime() (float64, error) {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0, err
	}
	fields := strings.Fields(string(data))
	if len(fields) == 0 {
		return 0, fmt.Errorf("unexpected /proc/uptime format")
	}
	return strconv.ParseFloat(fields[0], 64)
}

func formatUptime(seconds float64) string {
	total := int(seconds)
	days := total / 86400
	hours := (total % 86400) / 3600
	minutes := (total % 3600) / 60
	secs := total % 60
	return fmt.Sprintf("%dd %dh %dm %ds", days, hours, minutes, secs)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	uptimeSec, err := getUptime()
	if err != nil {
		http.Error(w, `{"error": "failed to read system uptime"}`, http.StatusInternalServerError)
		return
	}

	resp := StatusResponse{
		Hostname:      hostname,
		UptimeSeconds: uptimeSec,
		UptimeHuman:   formatUptime(uptimeSec),
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func main() {
	http.HandleFunc("/", statusHandler)
	http.HandleFunc("/health", statusHandler)

	port := "8000"
	log.Printf("Starting server on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
