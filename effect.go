package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
	"time"
)

func MemoryEffect(memorySize string, memoryStep string) {
	// 初始化
	fmt.Printf("\nMemory:\n")
	start := time.Now()

	// 转为字节数
	size, err := CalcByteNumber(memorySize)
	if err != nil {
		log.Fatalln(err)
	}
	step, err := CalcByteNumber(memoryStep)
	if err != nil {
		log.Fatalln(err)
	}

	// 分配内存
	data := make([]byte, size)
	for i := 0; i < size; i += step {
		data[i] = byte(1)
		seconds := time.Since(start).Seconds()
		percent := float32(i+1) / float32(size) * 100.0
		fmt.Printf("  Allocated: %d/%d bytes, Take time: %.1f second, Completed: %.1f%%\r",
			i+1,
			size,
			seconds,
			percent,
		)
	}

	// 最后一个字节分配内存
	data[size-1] = byte(1)
	seconds := time.Since(start).Seconds()

	// 输出信息
	fmt.Printf("                                                                                    \r")
	fmt.Printf("  Allocated: %d/%d bytes\n", size, size)
	fmt.Printf("  Take time: %.1f second\n", seconds)
}

func CpuEffect(cpuCore float64, validCore int) error {
	// 最多可使用逻辑CPU核心数
	runtime.GOMAXPROCS(validCore)

	// 使用率不可超过最大使用率
	if cpuCore > float64(validCore) {
		cpuCore = float64(validCore)
	}

	// CPU限制指标
	var (
		PeriodUs int
		QuotaUs  int
	)
	PeriodUs = 250000
	QuotaUs = int(float64(PeriodUs) * cpuCore)

	// 添加CPU使用率限制
	name := path.Join("effect", strconv.FormatInt(time.Now().Unix(), 10)) // 生成一个时间戳作为目录名
	err := cgroup.Add(
		Subsystem{
			Type:       "cpu",
			Name:       name,
			File:       "cpu.cfs_period_us",
			Values:     []string{strconv.Itoa(PeriodUs)},
			ProcessIds: []int32{int32(os.Getpid())},
		},
		Subsystem{
			Type:       "cpu",
			Name:       name,
			File:       "cpu.cfs_quota_us",
			Values:     []string{strconv.Itoa(QuotaUs)},
			ProcessIds: []int32{int32(os.Getpid())},
		},
	)

	if err != nil {
		return err
	}

	// 消耗CPU
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				time.Sleep(time.Nanosecond)
			}
		}()
	}

	// 输出信息
	fmt.Printf("\nCPU:\n")
	fmt.Printf("  Percent used : %d%%/%d%%\n", int(cpuCore*100), runtime.NumCPU()*100)
	fmt.Printf("  Validly cores: %d/%d\n", validCore, runtime.NumCPU())

	return nil
}

func DiskEffect(name string, diskSize string, bufferSize string) error {
	// 初始化
	fmt.Printf("\nDisk:\n")
	start := time.Now()

	// 计算字节数
	dsize, err := CalcByteNumber(diskSize)
	if err != nil {
		return err
	}
	bsize, err := CalcByteNumber(bufferSize)
	if err != nil {
		return nil
	}

	// 计算分块数量和剩余字节数
	chunk := dsize / bsize
	fsize := dsize % bsize

	// 打开文件
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	// 分块写入数据
	for i := 0; i < chunk; i++ {
		buffer := make([]byte, bsize)
		_, err = f.Write(buffer)
		if err != nil {
			return err
		}
		seconds := time.Since(start).Seconds()
		percent := float64(bsize*(i+1)) / float64(dsize) * 100.0
		fmt.Printf("  Allocated: %d/%d bytes, Take time: %.1f second, Completed: %.1f%%\r",
			bsize*(i+1),
			dsize,
			seconds,
			percent,
		)
	}

	// 写入剩余数据
	if fsize > 0 {
		buffer := make([]byte, fsize)
		_, err = f.Write(buffer)
		if err != nil {
			return err
		}
		seconds := time.Since(start).Seconds()
		fmt.Printf("  Allocated: %d/%d bytes, Take time: %.1f second, Completed: %.1f%%\r",
			dsize,
			dsize,
			seconds,
			100.0,
		)
	}

	return nil
}
