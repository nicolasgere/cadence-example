package main

import (
	t "cadence-example/adaptater"
	"cadence-example/config"
	"cadence-example/workflows/overdub"
	"fmt"
	"go.uber.org/cadence/worker"
	"go.uber.org/zap"
)

func startWorkers(h *t.CadenceAdapter, taskList string) {
	// Configure worker options.
	workerOptions := worker.Options{
		MetricsScope: h.Scope,
		Logger:       h.Logger,
	}

	cadenceWorker := worker.New(h.ServiceClient, h.Config.Domain, taskList, workerOptions)
	err := cadenceWorker.Start()
	if err != nil {
		h.Logger.Error("Failed to start workers.", zap.Error(err))
		panic("Failed to start workers")
	}
}

func main() {
	fmt.Println("Starting Worker..")
	var appConfig config.AppConfig
	appConfig.Setup()
	var cadenceClient t.CadenceAdapter
	cadenceClient.Setup(&appConfig.Cadence)

	startWorkers(&cadenceClient, overdub.TaskListName)
	// The workers are supposed to be long running process that should not exit.
	select {}
}
