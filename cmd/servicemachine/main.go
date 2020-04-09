package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

//-----------/dev/device

type Blockdevice struct {
	Devices []DeviceParent `json:"blockdevices"`
}

type DeviceParent struct {
	Name     string           `json:"name"`
	Type     string           `json:"type"`
	Size     string           `json:"size"`
	Children []DeviceChildren `json:"children"`
}

type DeviceChildren struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Size       string `json:"size"`
	Mountpoint string `json:"mountpoint"`
	Available  string `json:"available"`
	Used       string `json:"used"`
	Percent    string `json:"percent"`
}

type DiskStatus struct {
	All  uint64 `json:"all"`
	Used uint64 `json:"used"`
	Free uint64 `json:"free"`
}

//----------------------

//----cpu device

/*
The meanings of the columns are as follows, from left to right:
1st column : user = normal processes executing in user mode
2nd column : nice = niced processes executing in user mode
3rd column : system = processes executing in kernel mode
4th column : idle = twiddling thumbs
5th column : iowait = waiting for I/O to complete
6th column : irq = servicing interrupts
7th column : softirq = servicing softirqs
*/

type Blockcpu struct {
	CpuDevice []CpuDetails `json:"cpublock"`
}

type CpuDetails struct {
	Name    string `json:"name"`
	User    string `json:"user"`
	Nice    string `json:"nice"`
	System  string `json:"system"`
	Idle    string `json:"idle"`
	Iowait  string `json:"iowait"`
	Irq     string `json:"irq"`
	Softirq string `json:"softirq"`
}

// disk usage of path/disk
func DiskUsage(path string) (disk DiskStatus) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)
	if err != nil {
		return
	}
	disk.All = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.All - disk.Free
	return
}

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

func GetStorageDetails() Blockdevice {
	var blkdevice Blockdevice

	out, err := exec.Command("lsblk", "-I 8", "-J", "--output=name,type,size,mountpoint").Output()
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(out, &blkdevice)
	if err != nil {
		fmt.Println("An error occured: %v", err)
		os.Exit(1)
	}
	for indexParent, _ := range blkdevice.Devices {
		for indexChildren, _ := range blkdevice.Devices[indexParent].Children {
			if len(blkdevice.Devices[indexParent].Children[indexChildren].Mountpoint) > 0 {
				disk := DiskUsage(blkdevice.Devices[indexParent].Children[indexChildren].Mountpoint)
				blkdevice.Devices[indexParent].Children[indexChildren].Size = strconv.FormatFloat(float64(disk.All)/float64(GB), 'f', 2, 64) + "GB"
				blkdevice.Devices[indexParent].Children[indexChildren].Available = strconv.FormatFloat(float64(disk.Free)/float64(GB), 'f', 2, 64) + "GB"
				blkdevice.Devices[indexParent].Children[indexChildren].Used = strconv.FormatFloat(float64(disk.Used)/float64(GB), 'f', 2, 64) + "GB"
				blkdevice.Devices[indexParent].Children[indexChildren].Percent =
					strconv.FormatFloat(float64(disk.Used)/float64(disk.All)*100, 'f', 2, 64) + "%"
			}
		}
	}
	return blkdevice
}

func CalculateMediumCpu(info map[string][]float64) Blockcpu {
	//fmt.Println(info)
	var blockCpu Blockcpu
	total := make(map[string]float64)

	for k, v := range info {
		total[k] = 0
		for i := 0; i < len(v); i++ {
			total[k] += v[i]
		}
	}

	for k, v := range info {
		var tdetails CpuDetails
		tdetails.Name = k
		tdetails.User = strconv.FormatFloat(v[0]/total[k], 'f', 4, 64)
		tdetails.Nice = strconv.FormatFloat(v[1]/total[k], 'f', 4, 64)
		tdetails.System = strconv.FormatFloat(v[2]/total[k], 'f', 4, 64)
		tdetails.Idle = strconv.FormatFloat(v[3]/total[k], 'f', 4, 64)
		tdetails.Iowait = strconv.FormatFloat(v[4]/total[k], 'f', 4, 64)
		tdetails.Irq = strconv.FormatFloat(v[5]/total[k], 'f', 4, 64)
		tdetails.Softirq = strconv.FormatFloat(v[6]/total[k], 'f', 4, 64)
		blockCpu.CpuDevice = append(blockCpu.CpuDevice, tdetails)
	}
	// jsonA, _ := json.Marshal(blockCpu)
	// fmt.Println(string(jsonA))
	return blockCpu
}

func GetCpuDetails() Blockcpu {

	firstStat := GetProcStat()

	time.Sleep(time.Second * 3) //3 seconds like top

	secondStat := GetProcStat()

	fmt.Println(firstStat)
	fmt.Println(secondStat)
	difference := DifferenceProcStat(firstStat, secondStat)

	// jsonA, _ := json.Marshal(blockCpu)
	// fmt.Println(string(jsonA))
	return CalculateMediumCpu(difference)

}

func DifferenceProcStat(firstStat, secondStat string) map[string][]float64 {

	scannerFirst := bufio.NewScanner(strings.NewReader(firstStat))
	first := make(map[string][]float64)
	for scannerFirst.Scan() {
		slice := strings.Split(scannerFirst.Text(), " ")
		if strings.Contains(slice[0], "cpu") {
			// hack
			// there 2 spaces after regular cpu
			// cpu  177027 238 43905 18712690 5554 0 2065 0 0 0
			x := 0
			if slice[0] == "cpu" {
				x = 1
			}
			for i := 1 + x; i < len(slice)-x; i++ {

				floatValue, _ := strconv.ParseFloat(slice[i], 32)
				first[slice[0]] = append(first[slice[0]], floatValue)
			}
		}
	}
	second := make(map[string][]float64)
	scannerSecond := bufio.NewScanner(strings.NewReader(secondStat))
	for scannerSecond.Scan() {
		slice := strings.Split(scannerSecond.Text(), " ")
		if strings.Contains(slice[0], "cpu") {
			// hack
			// there 2 spaces after regular cpu
			// cpu  177027 238 43905 18712690 5554 0 2065 0 0 0
			x := 0
			if slice[0] == "cpu" {
				x = 1
			}
			for i := 1 + x; i < len(slice)-x; i++ {
				floatValue, _ := strconv.ParseFloat(slice[i], 32)
				second[slice[0]] = append(second[slice[0]], floatValue)
			}
		}
	}
	// fmt.Println(first)
	// fmt.Println(second)

	difference := make(map[string][]float64)

	for k, v := range first {
		for i := 0; i < len(v); i++ {
			difference[k] = append(difference[k], second[k][i]-v[i])
		}
	}
	return difference
}

func GetProcStat() string {
	out, err := exec.Command("cat", "/proc/stat").Output()
	if err != nil {
		log.Fatal(err)
	}

	stringOut := string(out)
	return stringOut

}

func main() {

	jsonA, _ := json.Marshal(GetCpuDetails())
	fmt.Println(string(jsonA))

	blockdevice := GetStorageDetails()
	b, _ := json.Marshal(blockdevice)
	fmt.Println(string(b))
}
