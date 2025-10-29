package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

var errorCount = 0

func main() {
	// Выполняем запрос и обрабатываем статистику
	fetchAndProcessStats()
}

func fetchAndProcessStats() {
	// Отправляем GET запрос к серверу
	response, err := http.Get("http://srv.msk01.gigacorp.local/_stats")
	if err != nil {
		fmt.Printf("Ошибка при выполнении запроса: %v\n", err)
		handleError()
		return
	}
	defer response.Body.Close()

	// Проверяем статус код
	if response.StatusCode != 200 {
		fmt.Printf("Ошибка: получен статус код %d вместо 200\n", response.StatusCode)
		handleError()
		return
	}

	// Читаем тело ответа
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Ошибка при чтении ответа: %v\n", err)
		handleError()
		return
	}

	// Преобразуем тело ответа в строку и разбиваем по запятым
	bodyStr := strings.TrimSpace(string(body))
	values := strings.Split(bodyStr, ",")
	
	// Проверяем, что получили достаточно значений
	if len(values) < 7 {
		handleError()
		return
	}

	// Сбрасываем счетчик ошибок при успешном запросе
	errorCount = 0

	// 1. Обрабатываем Load Average (первое значение - индекс 0)
	loadAvgStr := strings.TrimSpace(values[0])
	loadAvg, err := strconv.ParseFloat(loadAvgStr, 64)
	if err != nil {
		handleError()
		return
	}

	// Проверяем Load Average
	if loadAvg > 30 {
		fmt.Printf("Load Average is too high: %.2f\n", loadAvg)
	}

	// 2. Обрабатываем память (индексы 1 и 2)
	totalMemStr := strings.TrimSpace(values[1])
	totalMem, err := strconv.ParseUint(totalMemStr, 10, 64)
	if err != nil {
		handleError()
		return
	}

	usedMemStr := strings.TrimSpace(values[2])
	usedMem, err := strconv.ParseUint(usedMemStr, 10, 64)
	if err != nil {
		handleError()
		return
	}

	// Проверяем использование памяти
	if totalMem > 0 {
		memoryUsagePercent := float64(usedMem) / float64(totalMem) * 100
		if memoryUsagePercent > 80 {
			fmt.Printf("Memory usage too high: %.2f%%\n", memoryUsagePercent)
		}
	}

	// 3. Обрабатываем диск (индексы 3 и 4)
	totalDiskStr := strings.TrimSpace(values[3])
	totalDisk, err := strconv.ParseUint(totalDiskStr, 10, 64)
	if err != nil {
		handleError()
		return
	}

	usedDiskStr := strings.TrimSpace(values[4])
	usedDisk, err := strconv.ParseUint(usedDiskStr, 10, 64)
	if err != nil {
		handleError()
		return
	}

	// Проверяем использование диска
	if totalDisk > 0 {
		freeDisk := totalDisk - usedDisk
		freeDiskMB := freeDisk / (1024 * 1024) // Конвертируем в мегабайты
		diskUsagePercent := float64(usedDisk) / float64(totalDisk) * 100
		
		if diskUsagePercent > 90 {
			fmt.Printf("Free disk space is too low: %d Mb left\n", freeDiskMB)
		}
	}

	// 4. Обрабатываем сеть (индексы 5 и 6)
	totalNetStr := strings.TrimSpace(values[5])
	totalNet, err := strconv.ParseUint(totalNetStr, 10, 64)
	if err != nil {
		handleError()
		return
	}

	usedNetStr := strings.TrimSpace(values[6])
	usedNet, err := strconv.ParseUint(usedNetStr, 10, 64)
	if err != nil {
		handleError()
		return
	}

	// Проверяем загруженность сети
	if totalNet > 0 {
		freeNet := totalNet - usedNet
		freeNetMbit := float64(freeNet) / (1024 * 1024 / 8) // Конвертируем байты/сек в мегабиты/сек
		netUsagePercent := float64(usedNet) / float64(totalNet) * 100
		
		if netUsagePercent > 90 {
			fmt.Printf("Network bandwidth usage high: %.2f Mbit/s available\n", freeNetMbit)
		}
	}
}

func handleError() {
	errorCount++
	if errorCount >= 3 {
		fmt.Printf("Unable to fetch server statistic\n")
	}
}