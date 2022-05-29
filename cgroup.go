package main

import (
	"os"
	"path"
	"strconv"
	"strings"
)

type Cgroup struct {
	Path string                 // Cgroup目录
	m    map[string][]Subsystem // 子系统
}

type Subsystem struct {
	// 资源类型, ls -l /sys/fs/cgroup/下的目录便是资源类型
	// 支持输入 相对目录名(比如cpu) 或 绝对目录名(/sys/fs/cgroup/cpu)
	// 为了开发方便,不管输入啥都使用绝对目录名存储
	Type       string
	Name       string   // 子系统名称
	File       string   // 具体的文件名,比如 cpu.rt_period_us
	Values     []string // 文件中写入的值
	ProcessIds []int32  // 要限制的进程ID
}

// 添加子系统并使之生效
// 若指定ProcessIds则要求PID必须存在
func (c *Cgroup) Add(first Subsystem, others ...Subsystem) error {
	// 合并
	subs := append([]Subsystem{first}, others...)

	// 遍历
	for _, sys := range subs {
		// 构造目录
		if !strings.HasPrefix(sys.Type, "/") {
			sys.Type = path.Join(c.Path, sys.Type)
		}
		dir := path.Join(sys.Type, sys.Name)

		// 若目录不存在则创建
		if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
			err = os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return err
			}
		}

		// 目录成功创建之后记录下来
		c.m[sys.Name] = append(c.m[sys.Name], sys)

		// 写入文件
		file := path.Join(dir, sys.File)
		for _, value := range sys.Values {
			err := os.WriteFile(file, []byte(value), os.ModePerm)
			if err != nil {
				return err
			}
		}

		//获取进程的所有线程ID，并写入tasks文件
		for _, pid := range sys.ProcessIds {
			tids, err := GetThreadIds(pid)
			if err != nil {
				return err
			}
			for _, id := range tids {
				s := strconv.FormatInt(int64(id), 10)
				err := os.WriteFile(path.Join(dir, "tasks"), []byte(s), os.ModePerm)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// 删除指定条目
// 要求ProcessIds中的PID必须不存在,否则会删除失败
func (c *Cgroup) Delete(subname string) error {
	// 不存在直接返回
	v, ok := c.m[subname]
	if !ok {
		return nil
	}

	// 删除目录
	for _, sys := range v {
		dir := path.Join(sys.Type, sys.Name)
		err := os.Remove(dir)
		if err != nil {
			return err
		}
	}

	return nil
}

// 删除所有条目, 忽略错误
func (c *Cgroup) Clean() {
	for subname, _ := range c.m {
		_ = c.Delete(subname)
	}
}

func NewCgroup() *Cgroup {
	c := &Cgroup{
		Path: "/sys/fs/cgroup/",
		m:    make(map[string][]Subsystem),
	}
	return c
}

var cgroup = NewCgroup()
