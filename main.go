package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var errorCount int

func main() {
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

	load, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		countError()
		return
	}
	if load > 30 {
		fmt.Printf("Load Average is too high: %.0f\n", load) // ИСПРАВЛЕНО: %.0f вместо %.2f
	}

	totalMem, err1 := strconv.ParseUint(strings.TrimSpace(parts[1]), 10, 64)
	usedMem, err2 := strconv.ParseUint(strings.TrimSpace(parts[2]), 10, 64)
	if err1 != nil || err2 != nil || totalMem == 0 {
		countError()
		return
	}
	memUsage := float64(usedMem) / float64(totalMem) * 100
	if memUsage > 80 {
		fmt.Printf("Memory usage too high: %.0f%%\n", memUsage) // ИСПРАВЛЕНО: %.0f вместо %.2f
	}

	totalDisk, err1 := strconv.ParseUint(strings.TrimSpace(parts[3]), 10, 64)
	usedDisk, err2 := strconv.ParseUint(strings.TrimSpace(parts[4]), 10, 64)
	if err1 != nil || err2 != nil || totalDisk == 0 {
		countError()
		return
	}
	diskUsage := float64(usedDisk) / float64(totalDisk) * 100
	if diskUsage > 90 {
		freeDiskMB := (totalDisk - usedDisk) / (1024 * 1024)
		fmt.Printf("Free disk space is too low: %d Mb left\n", freeDiskMB)
	}

	totalNet, err1 := strconv.ParseUint(strings.TrimSpace(parts[5]), 10, 64)
	usedNet, err2 := strconv.ParseUint(strings.TrimSpace(parts[6]), 10, 64)
	if err1 != nil || err2 != nil || totalNet == 0 {
		countError()
		return
	}
	netUsage := float64(usedNet) / float64(totalNet) * 100
	if netUsage > 90 {
		freeNetMbit := float64(totalNet-usedNet) / (1024 * 1024 / 8)
		fmt.Printf("Network bandwidth usage high: %.0f Mbit/s available\n", freeNetMbit) // ИСПРАВЛЕНО: %.0f вместо %.2f
	}
}

func countError() {
	errorCount++
	if errorCount >= 3 {
		fmt.Printf("Unable to fetch server statistic\n")
	}
}