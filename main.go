package main

import (
	"os"
	"runtime/pprof"

	"runtime"

	"github.com/apex/log"

	"github.com/tmacychen/ContentSearch/cmd"
)

func main() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatalf("create cpu pprof error :%v\n", err)
		os.Exit(1)
	}
	defer f.Close()
	m, err := os.Create("mem.prof")
	if err != nil {
		log.Fatalf("create mem pprof error :%v\n", err)
		os.Exit(1)
	}
	defer m.Close()
	runtime.GC()
	if e := pprof.WriteHeapProfile(m); e != nil {
		log.Fatalf("mem prof err :%v\n", e)
	}
	pprof.StartCPUProfile(f)
	cmd.Execute()
	pprof.StopCPUProfile()
}
