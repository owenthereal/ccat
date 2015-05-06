package main

import (
	"log"
	"os"
	"runtime/pprof"

	_ "github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/dep"
	_ "github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/scan"
	"github.com/jingweno/ccat/Godeps/_workspace/src/sourcegraph.com/sourcegraph/srclib/src"
)

func main() {
	if cpuprof := os.Getenv("CPUPROF"); cpuprof != "" {
		f, err := os.Create(cpuprof)
		if err != nil {
			log.Fatal("CPUPROF:", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("StartCPUProfile:", err)
		}
		defer func() {
			pprof.StopCPUProfile()
			f.Close()
		}()
	}

	if err := src.Main(); err != nil {
		os.Exit(1)
	}
}
