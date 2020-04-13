package main

import (
	"encoding/json"
	"fmt"
)

func main() {

	// jsonA, _ := json.Marshal(GetCpuDetails())
	// fmt.Println(string(jsonA))

	// blockdevice := GetStorageDetails()
	// b, _ := json.Marshal(blockdevice)
	// fmt.Println(string(b))

	var blockMemory BlockMemory
	blockMemory.GetMemoryDetails()
	memory, _ := json.Marshal(blockMemory)
	fmt.Println(string(memory))
}
