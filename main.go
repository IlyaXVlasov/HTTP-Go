package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var errorCount int

func main() {
	for i := 0; i < 10; i++ {
		fetchAndProcessStats()
		time.Sleep(100 * time.Millisecond)
	}
}


func fetchAndProcessStats() {
	response, err := http.Get("http://srv.msk01.gigacorp.local/_stats")
	if err != nil {
		countError()
		return
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		countError()
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		countError()
		return
	}

	processData(string(body))
}

func processData(data string) {
	parts := strings.Split(strings.TrimSpace(data), ",")
	if len(parts) < 7 {
		countError()
		return
	}

	errorCount = 0

	// Load Average
	load, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err == nil && load > 30 {
		fmt.Printf("Load Average is too high: %.0f\n", load)
	}

	// Memory
	totalMem, err1 := strconv.ParseUint(strings.TrimSpace(parts[1]), 10, 64)
	usedMem, err2 := strconv.ParseUint(strings.TrimSpace(parts[2]), 10, 64)
	if err1 == nil && err2 == nil && totalMem > 0 {
		memUsagePercent := (usedMem * 100) / totalMem
		if memUsagePercent > 80 {
			fmt.Printf("Memory usage too high: %d%%\n", memUsagePercent)
		}
	}

	// Disk
	totalDisk, err1 := strconv.ParseUint(strings.TrimSpace(parts[3]), 10, 64)
	usedDisk, err2 := strconv.ParseUint(strings.TrimSpace(parts[4]), 10, 64)
	if err1 == nil && err2 == nil && totalDisk > 0 {
		diskUsagePercent := (usedDisk * 100) / totalDisk
		if diskUsagePercent > 90 {
			freeDiskBytes := totalDisk - usedDisk
			freeDiskMB := freeDiskBytes / (1024 * 1024)
			fmt.Printf("Free disk space is too low: %d Mb left\n", freeDiskMB)
		}
	}


// Network
totalNet, err1 := strconv.ParseUint(strings.TrimSpace(parts[5]), 10, 64)
usedNet, err2 := strconv.ParseUint(strings.TrimSpace(parts[6]), 10, 64)

if err1 == nil && err2 == nil && totalNet > 0 {
    netUsagePercent := (usedNet * 100) / totalNet
    if netUsagePercent > 90 {
        freeNetBytes := totalNet - usedNet
        
        // Округление до целого
        freeNetMbit := (freeNetBytes * 8 + 500000) / 1000000
        
        fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", freeNetMbit)
    }
}
}

func countError() {
    errorCount++
    if errorCount >= 3 {
        fmt.Printf("Unable to fetch server statistic\n")
    }
}