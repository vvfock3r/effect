# effect [![Go Report Card](https://goreportcard.com/badge/github.com/vvfock3r/effect)](https://goreportcard.com/report/github.com/vvfock3r/effect)

Effect是一个用于自定义系统资源利用率的工具，编写这个工具的目的在于学习Go、Linux、Docker


## Install

Go：requires Go 1.18 or later

```bash
go install github.com/vvfock3r/effect@latest
```

> 原理：
>
> （1）下载代码到`GOPATH/pkg/mod/github.com/vvfock3r/`目录下
>
> （2）并编译代码，将二进制文件到`GOPATH/bin`目录下
>
> （3）查看GOPATH目录方法：`go env GOPATH`



## Usage

```bash
[root@localhost ~]# ./effect -h
Usage: ./effect OPTIONS 

Effect is a tool for customizing system resource utilization.
For details, refer to https://github.com/vvfock3r/effect

OPTIONS:
  -cpu float
        CPU utilization
  -cpu-spread int
        Number of CPU propagation (default 4)
  -memory string
        Allocated memory size
  -memory-stride string
        Memory increase stride (default "1024kb")

Example:
  ./effect --memory 2g --cpu 1.5
```

### Memory

注意事项：

* 直接向系统申请内存，所以此功能支持Windows、Linux等
* 程序目前无法区分交换分区和物理内存，所以这里是包含物理内存和交换分区的，在测试时最好将交换分区关掉，以便能更好的观测到结果
* 程序内部以1024为转换单位，若指定为2g实际为2048mb；而常用的top命令默认以1000为转换单位
* 只支持整数单位，不支持比如1.5g，请使用1500mb或1536mb

<br />

（1）向系统申请`2G`内存（默认每次申请`1024kb`）

![image-20220529131502995](https://tuchuang-1257805459.cos.ap-shanghai.myqcloud.com/image-20220529131502995.png)

（2）向系统申请`2G`内存，每次申请`4kb`（可以看到速度有明显下降）

![image-20220529131739126](https://tuchuang-1257805459.cos.ap-shanghai.myqcloud.com/image-20220529131739126.png)

### CPU

注意事项：

* 程序内部会将CPU跑满，然后再通过Linux Cgroup来限制CPU使用率，所以该功能不支持Windows

<br />

（1）申请1.5核的CPU（默认会平均分布到各个核心上）

![image-20220529132309015](https://tuchuang-1257805459.cos.ap-shanghai.myqcloud.com/image-20220529132309015.png)

（2）申请1.5核的CPU（仅分布在两个核心上）

![image-20220529132550375](https://tuchuang-1257805459.cos.ap-shanghai.myqcloud.com/image-20220529132550375.png)





## ToDo

* 添加磁盘使用率
* 增加百分比数据格式
* 增加平均负载使用率
  * CPU引起负载升高
  * 短时应用引起负载升高（top中是看不到进程名的）
  * IO等待引起负载升高
* other ...
