package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
	"time"
)

// 一段有问题的代码
func logicCode() {
	var c chan int
	for {
		select {
		case v := <-c:
			fmt.Printf("recv from chan, value:%v\n", v)
		default:
			//time.Sleep(2 * time.Second)
		}
	}
}

func main() {
	// 两个标志位: 是否开启CPU和内存的标志位
	var isCPUPprof bool
	var isMemPprof bool

	// 命令行参数定义
	flag.BoolVar(&isCPUPprof, "cpu", false, "turn cpu pprof on")
	flag.BoolVar(&isMemPprof, "mem", false, "turn mem pprof on")
	flag.Parse()

	// 是否开启CPUprofile
	if isCPUPprof {
		// 在当前路径建立一个文件
		file, err := os.Create("./cpu.pprof")
		if err != nil {
			fmt.Printf("create cpu pprof failed, err:%v\n", err)
			return
		}
		// 往文件中记录CPU proofile信息
		pprof.StartCPUProfile(file)
		defer func() {
			pprof.StopCPUProfile()
			file.Close()
		}()
	}
	for i := 0; i < 8; i++ {
		go logicCode()
	}

	// 程序跑20s
	time.Sleep(20 * time.Second)

	// 是否开启内存profile
	if isMemPprof {
		file, err := os.Create("./mem.pprof")
		if err != nil {
			fmt.Printf("create mem pprof failed, err:%v\n", err)
			return
		}
		pprof.WriteHeapProfile(file)
		file.Close()
	}

}
