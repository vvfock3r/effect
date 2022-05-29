package main

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/process"
	"regexp"
	"strconv"
	"strings"
)

// 根据输入(比如2g、10kb等)来计算字节数
func CalcByteNumber(s string) (ret int, err error) {
	// 转为大写
	input := strings.ToUpper(s)

	// 提取数值和单位
	number := regexp.MustCompile("([0-9]+)")
	unit := regexp.MustCompile("([A-Z]+)")
	nString := strings.TrimSpace(number.FindString(input))
	uString := strings.TrimSpace(unit.FindString(input))
	if len(nString) == 0 || len(uString) == 0 {
		return ret, fmt.Errorf("CalcByteNumber error: %s\n", s)
	}

	// 字符串数值转为数字类型
	n, err := strconv.Atoi(nString)
	if err != nil {
		return ret, fmt.Errorf("CalcByteNumber error: %s\n", s)
	}

	// 单位转换
	switch uString {
	case "K", "KB", "KIB":
		ret = n * 1024
	case "M", "MB", "MIB":
		ret = n * 1024 * 1024
	case "G", "GB", "GIB":
		ret = n * 1024 * 1024 * 1024
	default:
		err = fmt.Errorf("CalcByteNumber error: %s\n", s)
	}

	return ret, err
}

// 获取进程所有线程ID, 等同于 top -H -p<pid>，不支持Windows
func GetThreadIds(pid int32) (tids []int32, err error) {
	// 实例化进程对象
	p, err := process.NewProcess(pid)
	if err != nil {
		return nil, fmt.Errorf("GetThreadIds error: %s\n", err.Error())
	}

	// 获取所有线程
	threads, err := p.Threads()
	if err != nil {
		return nil, fmt.Errorf("GetThreadIds error: %s\n", err.Error())
	}

	// 提取所有线程ID
	for id, _ := range threads {
		tids = append(tids, id)
	}

	return tids, err
}
