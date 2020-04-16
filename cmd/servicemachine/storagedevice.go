package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"syscall"
)

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
					strconv.FormatFloat(float64(disk.Used)/float64(disk.All)*100, 'f', 2, 64)
			}
		}
	}
	return blkdevice
}
