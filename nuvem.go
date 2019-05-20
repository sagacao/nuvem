package main

import (
	"nuvem/engine/logger"

	//"QGameServer/network"
	"fmt"
	_ "nuvem/engine/utils"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
)

func cpuProfile() {
	f, err := os.OpenFile("cpu.prof", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		logger.Fatal(err)
	}
	defer f.Close()

	logger.Debug("CPU Profile started")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	//time.Sleep(60 * time.Second)
	//logger.Debug("CPU Profile stopped")
}

func heapProfile() {
	f, err := os.OpenFile("heap.prof", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		logger.Fatal(err)
	}
	defer f.Close()

	//time.Sleep(30 * time.Second)
	pprof.WriteHeapProfile(f)
	logger.Debug("Heap Profile generated")
}

func saveHeapProfile() {
	runtime.GC()
	f, err := os.Create(fmt.Sprintf("prof/heap_%s.prof", time.Now().Format("2006_01_02_03_04_05")))
	if err != nil {
		return
	}
	defer f.Close()
	pprof.Lookup("heap").WriteTo(f, 1)
}

func main() {
	// runtime.GOMAXPROCS(runtime.NumCPU())

	// logger.SetRollingDaily(settings.ServerConfig.Env.ACCESS_LOG_PATH, settings.ServerConfig.Env.ACCESS_LOG_NAME)
	// if settings.ServerConfig.Env.DEBUG {
	// 	gin.SetMode(gin.DebugMode)
	// 	logger.SetConsole(true)
	// 	logger.SetLevel(logger.DEBUG)

	// 	//cpuProfile()
	// 	//heapProfile()
	// } else {
	// 	gin.SetMode(gin.ReleaseMode)
	// 	logger.SetLevel(logger.INFO)
	// }

	//network.RunConfigServer()
	//network.RunGMServer()
	//network.RunGate()
	//network.RunTCPServer()
}
