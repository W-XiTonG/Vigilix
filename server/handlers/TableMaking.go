package handlers

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"log"
	"path/filepath"
	"server/util"
)

func TableMaking(f *excelize.File, tablePaths, deliverTime, name string, data *util.ClientSystemMetrics,
	clientMsgID int, diskPartition util.DiskPartition) {
	// 创建一个工作表（默认会有一个名为 "Sheet1" 的工作表）
	index, _ := f.NewSheet("Sheet1")

	// 定义表头
	headers := []string{
		"Agent描述", "主机名", "操作系统类型", "系统架构", "内核版本",
		"CPU物理核心数", "CPU逻辑核心数", "(%)CPU总利用率", "(%)CPU用户态时间占比", "(%)CPU系统态时间占比",
		"(%)CPU空闲时间占比", "(MB)Memory内存总量", "(MB)Memory已用内存", "(%)Memory内存使用率",
		"Disk分区挂载点", "(MB)Disk分区空间", "(MB)Disk分区已用空间", "(%)Disk分区使用率",
		"(MB)Network总发送字节数", "(MB)Network总接收字节数",
		"(Load)1分钟平均负载", "(Load)5分钟平均负载", "(Load)15分钟平均负载",
	}
	integer, err := f.NewStyle(&excelize.Style{
		NumFmt: 0,
	})
	if err != nil {
		log.Fatalf("Error: 创建样式时出错: %v", err)
	}
	percentage, err := f.NewStyle(&excelize.Style{
		NumFmt: 2,
	})
	if err != nil {
		log.Fatalf("Error: 创建样式时出错: %v", err)
	}

	clientData := []any{
		name, data.Host.Hostname, data.Host.OS, data.Host.Kernel, data.Host.Architecture,
		data.CPU.PhysicalCores, data.CPU.LogicalCores, data.CPU.TotalUsage, data.CPU.UserMode, data.CPU.SystemMode, data.CPU.Idle,
		util.BytesToMB(data.Memory.Total), util.BytesToMB(data.Memory.Used), data.Memory.UsedPercent,
		diskPartition.MountPoint, util.BytesToMB(diskPartition.Total), util.BytesToMB(diskPartition.Used), diskPartition.UsedPercent,
		util.BytesToMB(data.Network.SentTotal), util.BytesToMB(data.Network.RecvTotal),
		data.SystemLoad.Load1, data.SystemLoad.Load5, data.SystemLoad.Load15,
	}

	// 设置表头单元格的值
	for col, header := range headers {
		cell, err := excelize.CoordinatesToCellName(col+1, 1)
		if err != nil {
			log.Printf("Error: %v", err)
			continue
		}
		if err := f.SetCellValue("Sheet1", cell, header); err != nil {
			log.Println(err)
		}
		cellData, err := excelize.CoordinatesToCellName(col+1, clientMsgID+1)
		if err := f.SetCellValue("Sheet1", cellData, clientData[col]); err != nil {
			log.Printf("Error: %v", err)
		}

		// 对特定列应用数字样式
		switch headers[col] {
		case "(%)CPU总利用率", "(%)CPU用户态时间占比", "(%)CPU系统态时间占比",
			"(%)CPU空闲时间占比", "(%)Memory内存使用率", "(%)Disk分区使用率",
			"(Load)1分钟平均负载", "(Load)5分钟平均负载", "(Load)15分钟平均负载":
			err = f.SetCellStyle("Sheet1", cellData, cellData, percentage)
			if err != nil {
				log.Printf("Error: %v", err)
			}
		case "CPU物理核心数", "CPU逻辑核心数":
			err = f.SetCellStyle("Sheet1", cellData, cellData, integer)
			if err != nil {
				log.Printf("Error: %v", err)
			}
		default:
			if util.IsNumber(clientData[col]) {
				err = f.SetCellStyle("Sheet1", cellData, cellData, percentage)
				if err != nil {
					log.Printf("Error: %v", err)
				}
			}
		}
	}
	// 设置活动工作表（默认显示的工作表）
	fileName := fmt.Sprintf("systemMetrics_%s.xlsx", deliverTime)
	f.SetActiveSheet(index)
	tableName := filepath.Join(tablePaths, fileName)
	// 保存文件
	if err := f.SaveAs(tableName); err != nil {
		log.Println("Error: Save failed error:", err)
	} else {
		log.Printf("The save was successful：%s", tableName)
		f = nil
	}
}
