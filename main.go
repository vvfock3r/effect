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
	memory       string
	memoryStride string
	cpu          float64 // 所有CPU核心总共使用率
	cpuSpread    int     // 最大可用逻辑CPU核心数
)

func init() {
	// 绑定内存参数
	flag.StringVar(&memory, "memory", "", "Allocated memory size")
	flag.StringVar(&memoryStride, "memory-stride", "1024kb", "Memory increase stride")

	// 绑定CPU参数
	flag.Float64Var(&cpu, "cpu", 0, "CPU utilization")
	flag.IntVar(&cpuSpread, "cpu-spread", runtime.NumCPU(), "Number of CPU propagation")

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
	if memory != "" {
		MemoryEffect(memory, memoryStride)
	}

	// CPU
	if cpu != 0.0 {
		err := CpuEffect(cpu, cpuSpread)
		if err != nil {
			log.Fatalln(err)
		}
	}

	// 提示信息
	fmt.Printf("\nPress <Ctrl-C> to quit.\n")

	for {
		time.Sleep(time.Second)
	}
}
