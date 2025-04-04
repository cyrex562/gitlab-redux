package service

import (
	"context"
	"strings"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

// SystemInfoService handles system information operations
type SystemInfoService struct {
	// TODO: Add necessary dependencies
}

// NewSystemInfoService creates a new instance of SystemInfoService
func NewSystemInfoService() *SystemInfoService {
	return &SystemInfoService{}
}

// SystemInfo represents the system information data
type SystemInfo struct {
	CPUs   []CPUInfo   `json:"cpus,omitempty"`
	Memory *MemoryInfo `json:"memory,omitempty"`
	Disks  []DiskInfo  `json:"disks"`
}

// CPUInfo represents CPU information
type CPUInfo struct {
	User   float64 `json:"user"`
	System float64 `json:"system"`
	Idle   float64 `json:"idle"`
	// Add other CPU-related fields as needed
}

// MemoryInfo represents memory information
type MemoryInfo struct {
	Total     uint64  `json:"total"`
	Used      uint64  `json:"used"`
	Free      uint64  `json:"free"`
	UsedPercent float64 `json:"used_percent"`
}

// DiskInfo represents disk information
type DiskInfo struct {
	BytesTotal uint64 `json:"bytes_total"`
	BytesUsed  uint64 `json:"bytes_used"`
	DiskName   string `json:"disk_name"`
	MountPath  string `json:"mount_path"`
}

var (
	excludedMountOptions = []string{
		"nobrowse",
		"read-only",
		"ro",
	}

	excludedMountTypes = []string{
		"autofs",
		"binfmt_misc",
		"bpf",
		"cgroup",
		"cgroup2",
		"configfs",
		"debugfs",
		"devfs",
		"devpts",
		"devtmpfs",
		"efivarfs",
		"fuse.gvfsd-fuse",
		"fuseblk",
		"fusectl",
		"hugetlbfs",
		"mqueue",
		"proc",
		"pstore",
		"rpc_pipefs",
		"securityfs",
		"sysfs",
		"tmpfs",
		"tracefs",
		"vfat",
	}
)

// GetSystemInfo retrieves system information including CPU, memory, and disk usage
func (s *SystemInfoService) GetSystemInfo(ctx context.Context) (*SystemInfo, error) {
	info := &SystemInfo{}

	// Get CPU information
	if cpuInfo, err := s.getCPUInfo(ctx); err == nil {
		info.CPUs = cpuInfo
	}

	// Get memory information
	if memInfo, err := s.getMemoryInfo(ctx); err == nil {
		info.Memory = memInfo
	}

	// Get disk information
	diskInfo, err := s.getDiskInfo(ctx)
	if err != nil {
		return nil, err
	}
	info.Disks = diskInfo

	return info, nil
}

// getCPUInfo retrieves CPU information
func (s *SystemInfoService) getCPUInfo(ctx context.Context) ([]CPUInfo, error) {
	percentages, err := cpu.Percent(0, false)
	if err != nil {
		return nil, err
	}

	var cpus []CPUInfo
	for _, p := range percentages {
		cpus = append(cpus, CPUInfo{
			User:   p.User,
			System: p.System,
			Idle:   p.Idle,
		})
	}
	return cpus, nil
}

// getMemoryInfo retrieves memory information
func (s *SystemInfoService) getMemoryInfo(ctx context.Context) (*MemoryInfo, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	return &MemoryInfo{
		Total:      v.Total,
		Used:       v.Used,
		Free:       v.Free,
		UsedPercent: v.UsedPercent,
	}, nil
}

// getDiskInfo retrieves disk information
func (s *SystemInfoService) getDiskInfo(ctx context.Context) ([]DiskInfo, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var disks []DiskInfo
	for _, partition := range partitions {
		// Skip excluded mount types
		if contains(excludedMountTypes, partition.Fstype) {
			continue
		}

		// Skip excluded mount options
		options := strings.Split(partition.Opts, ",")
		if hasExcludedOption(options) {
			continue
		}

		// Get disk usage information
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue // Skip if we can't get usage information
		}

		disks = append(disks, DiskInfo{
			BytesTotal: usage.Total,
			BytesUsed:  usage.Used,
			DiskName:   partition.Device,
			MountPath:  usage.Path,
		})
	}

	return disks, nil
}

// contains checks if a string slice contains a specific string
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

// hasExcludedOption checks if any of the mount options are in the excluded list
func hasExcludedOption(options []string) bool {
	for _, option := range options {
		if contains(excludedMountOptions, option) {
			return true
		}
	}
	return false
}
