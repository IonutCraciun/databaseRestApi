package main

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

/*
MemTotal:        8167640 kB   8167640 * 10^3
MemFree:         3071496 kB
MemAvailable:    6069264 kB
Buffers:           79844 kB
Cached:          2835212 kB

*/

const ratio = 1.024000

type BlockMemory struct {
	MemoryDeviceBlock MemoryDevice `json:"memoryblock"`
}

type MemoryDevice struct {
	MemTotal     string `json:"MemTotal"`
	MemFree      string `json:"MemFree"`
	MemAvailable string `json:"MemAvailable"`
	Buffers      string `json:"Buffers"`
	Cached       string `json:"Cached"`
	SwapTotal    string `json:"SwapTotal"`
	SwapFree     string `json:"SwapFree"`
	PercentUsage string `json:"PercentUsage"`
}

func (blockMemory *BlockMemory) GetMemoryDetails() {
	file := GetMemoryProc()
	var memoryDevice MemoryDevice
	scannerFile := bufio.NewScanner(strings.NewReader(file))
	for scannerFile.Scan() {
		slice := strings.Split(scannerFile.Text(), ":")

		switch strings.TrimSpace(slice[0]) {
		case "MemTotal":
			size, _ := strconv.ParseFloat(strings.Split(strings.TrimSpace(slice[1]), " ")[0], 32)
			memoryDevice.MemTotal = fmt.Sprintf("%.0f MB", (size*ratio*KB)/MB)
		case "MemFree":
			size, _ := strconv.ParseFloat(strings.Split(strings.TrimSpace(slice[1]), " ")[0], 32)
			memoryDevice.MemFree = fmt.Sprintf("%.0f MB", (size*ratio*KB)/MB)
		case "MemAvailable":
			size, _ := strconv.ParseFloat(strings.Split(strings.TrimSpace(slice[1]), " ")[0], 32)
			memoryDevice.MemAvailable = fmt.Sprintf("%.0f MB", (size*ratio*KB)/MB)
		case "Buffers":
			size, _ := strconv.ParseFloat(strings.Split(strings.TrimSpace(slice[1]), " ")[0], 32)
			memoryDevice.Buffers = fmt.Sprintf("%.0f MB", (size*ratio*KB)/MB)
		case "Cached":
			size, _ := strconv.ParseFloat(strings.Split(strings.TrimSpace(slice[1]), " ")[0], 32)
			memoryDevice.Cached = fmt.Sprintf("%.0f MB", (size*ratio*KB)/MB)
		case "SwapTotal":
			size, _ := strconv.ParseFloat(strings.Split(strings.TrimSpace(slice[1]), " ")[0], 32)
			memoryDevice.SwapTotal = fmt.Sprintf("%.0f MB", (size*ratio*KB)/MB)
		case "SwapFree":
			size, _ := strconv.ParseFloat(strings.Split(strings.TrimSpace(slice[1]), " ")[0], 32)
			memoryDevice.SwapFree = fmt.Sprintf("%.0f MB", (size*ratio*KB)/MB)
		}
	}

	blockMemory.MemoryDeviceBlock = memoryDevice
}

func GetMemoryProc() string {
	out, err := exec.Command("cat", "/proc/meminfo").Output()
	if err != nil {
		log.Fatal(err)
	}
	stringOut := string(out)
	return stringOut
}
