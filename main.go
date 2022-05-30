package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

var (
	memory       string // 总共申请内存大小
	memoryStride string // 每次申请内存大小

	cpu       float64 // 所有CPU核心总共使用率
	cpuSpread int     // 最大可用逻辑CPU核心数

	disk         string // 总共申请磁盘大小
	diskbuffer   string // 每次申请磁盘大小
	diskFileName string // 写入磁盘的文件名
)

func init() {
	// 绑定内存参数
	flag.StringVar(&memory, "memory", "", "Allocat memory size")
	flag.StringVar(&memoryStride, "memory-stride", "1024kb", "Memory increase stride")

	// 绑定CPU参数
	flag.Float64Var(&cpu, "cpu", 0, "CPU utilization")
	flag.IntVar(&cpuSpread, "cpu-spread", runtime.NumCPU(), "Number of CPU propagation")

	// 绑定磁盘参数
	flag.StringVar(&disk, "disk", "", "Allocat disk size")
	flag.StringVar(&diskbuffer, "disk-buffer", "1024kb", "Buffer size when writing file")
	flag.StringVar(&diskFileName, "disk-file", "effect.test.data", "File name")

	// Usage
	flag.Usage = usage
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s OPTIONS \n\n", os.Args[0])

	fmt.Fprintln(os.Stderr, "Effect is a tool for customizing system resource utilization.")
	fmt.Fprintln(os.Stderr, "For details, refer to https://github.com/vvfock3r/effect")
	fmt.Fprintln(os.Stderr, "")

	fmt.Fprintln(os.Stderr, "OPTIONS:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "")

	fmt.Fprintln(os.Stderr, "Example:")
	fmt.Fprintln(os.Stderr, "  ./effect --memory 2g --cpu 1.5")
}

func main() {
	// 选项解析
	flag.Parse()

	// 无选项
	if len(os.Args) <= 1 {
		flag.Usage()
		os.Exit(2)
	}

	// 含有未知的子命令时
	if len(flag.Args()) >= 1 {
		flag.Usage()
		os.Exit(2)
	}

	// 内存
	if len(memory) > 0 {
		MemoryEffect(memory, memoryStride)
	}

	// CPU
	if cpu != 0.0 {
		err := CpuEffect(cpu, cpuSpread)
		if err != nil {
			fmt.Println()
			log.Fatalln(err)
		}
	}

	// 磁盘
	if len(disk) > 0 {
		err := DiskEffect(diskFileName, disk, diskbuffer)
		if err != nil {
			fmt.Println()
			log.Fatalln(err)
		}
	}

	// 提示信息
	fmt.Printf("\nPress <Ctrl-C> to quit.\n")

	for {
		time.Sleep(time.Second)
	}
}
