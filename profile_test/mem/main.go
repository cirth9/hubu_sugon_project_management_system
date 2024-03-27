package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
	"sync"
)

var (
	cpu string
	mem string
)

func init() {
	flag.StringVar(&cpu, "cpu", "", "write cpu profile to file")
	flag.StringVar(&mem, "mem", "", "write mem profile to file")
}

func main() {
	flag.Parse()

	//采样 CPU 运行状态
	if cpu != "" {
		f, err := os.Create(cpu)
		if err != nil {
			log.Fatal(err)
		}
		_ = pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go workOnce(&wg)
	}
	wg.Wait()

	//采样内存状态
	if mem != "" {
		f, err := os.Create(mem)
		if err != nil {
			log.Fatal(err)
		}
		_ = pprof.WriteHeapProfile(f)
		f.Close()
	}
}

//func counter() {
//	var slice [100000]int
//	var c int
//	for i := 0; i < 100000; i++ {
//		c = i + 1 + 2 + 3 + 4 + 5
//		slice[i] = c
//	}
//	_ = slice
//}

func counter() {
	slice := make([]int, 0)
	var c int
	for i := 0; i < 100000; i++ {
		c = i + 1 + 2 + 3 + 4 + 5
		slice = append(slice, c)
	}
	_ = slice
}

func workOnce(wg *sync.WaitGroup) {
	counter()
	wg.Done()
}
