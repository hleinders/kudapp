package main

import (
	"fmt"
	"math/big"
	"math/rand"
	"time"
)

var (
	chStop               = make(chan bool, globalGFMaxCount)
	chWatch              = make(chan bool, 2)
	chResult             = make(chan bool, 2048)
	rndInit              = rand.NewSource(int64(time.Now().Nanosecond()))
	rnd                  = *rand.New(rndInit)
	localWorkerRunning   = false
	localWatchDogRunning = false
	localCollectRunning  = false
)

func fibonacci(id int) {
	// modified package example
	var limit big.Int
	var digits int64
	a := big.NewInt(0)
	b := big.NewInt(1)

	digits = int64(300000 + rnd.Intn(199999))
	limit.Exp(big.NewInt(10), big.NewInt(digits), nil)

	for a.Cmp(&limit) < 0 {
		a.Add(a, b)
		a, b = b, a
	}

	prVerboseInfo("Workload #%d: Finished Fibonacci Number (digits: %d)", id, digits)
	// prDebug("Workload #%d: Digits: %d\n Fibonacci No.: %s", id, digits, a)
}

func workload(stop chan (bool), id int) {
	prVerboseInfo("Workload #%d: started up!", id)

Loop:
	for {
		select {
		case <-stop:
			prVerboseInfo("Workload #%d: received stop signal", id)
			break Loop
		default:
			prVerboseInfo("%d: starting workload...", id)
			fibonacci(id)

			// add result
			chResult <- true

			// time.Sleep(10 * time.Second)
		}
	}

	prVerboseInfo("Workload #%d: stopped", id)
}

func gofuncsStartUp(chn chan bool, num int) error {
	if localWorkerRunning {
		prDebug("gofuncsStartAll: workload already running")
		return fmt.Errorf("workload already running")
	}

	for i := 0; i < num; i++ {
		k := i + 1
		go workload(chn, k)
	}
	localWorkerRunning = true
	return nil
}

func gofuncsStopAll(chn chan bool, num int) error {
	if !localWorkerRunning {
		prDebug("gofuncsStopAll: workload already stopped")
		return fmt.Errorf("workload already stopped")
	}

	for i := 0; i < num; i++ {
		chn <- true
	}
	localWorkerRunning = false
	return nil
}

func gofuncsWatchDog(quit chan bool, num int) {
	if localWatchDogRunning {
		prDebug("gofuncsWatchDog: watchdog already running")
	}

	startUpTime := time.Now()
	localWatchDogRunning = true
	prVerboseInfo("Worker watchdog: started up!")

Loop:
	// wait for globalGFMaxRuntime minutes and then kill workers
	for {
		select {
		case <-quit:
			prVerboseInfo("Worker watchdog: received stop signal")
			break Loop
		default:
			time.Sleep(3 * time.Second)
			duration := time.Since(startUpTime)
			if duration >= time.Duration(globalGFMaxRuntime*int(time.Minute)) {
				prVerboseInfo("ALERT: Worker watchdog: maximum runtime reached! stop all workers")
				displayErr(workoutStop(num))
				return
			}
			prDebug("watchdog still waiting...")
		}
	}
	localWatchDogRunning = false
	prVerboseInfo("Worker watchdog: stopped")
}

func gofuncsCollectResults() {
	if localCollectRunning {
		prDebug("gofuncsCollectResults: collector already running")
	}

	localCollectRunning = true
	prVerboseInfo("Result collector: started up!")

	// give workers time to start:
	time.Sleep(5 * time.Second)

Loop:
	for {
		select {
		case <-chResult:
			globalWorkerResult += 1
		default:
			if !localWorkerRunning {
				break Loop
			}
			time.Sleep(3 * time.Second)
		}
	}
	localCollectRunning = false
	prVerboseInfo("Result collector: stopped")
}

func workoutStart(num int) error {
	prVerboseInfo("Starting workout with %d worker(s)", num)

	// activate watchdocg
	go gofuncsWatchDog(chWatch, num)

	// activate results collector
	go gofuncsCollectResults()
	return gofuncsStartUp(chStop, num)
}

func workoutStop(num int) error {
	prVerboseInfo("Stopping workout with %d worker(s)", num)
	// deactivate watchdocg
	if localWatchDogRunning {
		chWatch <- true
	}

	return gofuncsStopAll(chStop, num)
}
