package util

import "errors"

// FindDiskPartition 函数用于从 []DiskPartition 中查找指定 MountPoint 的分区
func FindDiskPartition(partitions []DiskPartition, diskPartition string) (DiskPartition, bool) {
	for _, partition := range partitions {
		if partition.MountPoint == diskPartition {
			return partition, true
		}
	}
	return DiskPartition{}, false
}

// CheckAndAdd 函数用于检查整数是否在切片中，若不在则添加到切片
func CheckAndAdd(num string, nums []string) ([]string, error) {
	for _, v := range nums {
		if v == num {
			return nums, errors.New("already exists")
		}
	}
	return append(nums, num), nil
}

// RemoveElementByValue 根据元素值删除切片中的元素
func RemoveElementByValue(slice []string, value string) []string {
	result := make([]string, 0, len(slice))
	for _, v := range slice {
		if v != value {
			result = append(result, v)
		}
	}
	return result
}
