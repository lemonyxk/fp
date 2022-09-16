//go:build !windows
// +build !windows

/**
* @program: fp
*
* @description:
*
* @author: lemo
*
* @create: 2022-09-08 20:22
**/

package main

import (
	"context"
	"strconv"
	"syscall"

	"github.com/lemonyxk/console"
)

func kill(processes Processes) {

	var killValue = getArgs([]string{"-k", "--kill"})

	if killValue == "" {
		killValue = "15"
	}

	var killSignal, err = strconv.Atoi(killValue)
	if err != nil {
		console.Info(err)
		return
	}

	for i := 0; i < len(processes); i++ {
		var err = processes[i].process.SendSignalWithContext(context.Background(), syscall.Signal(killSignal))
		if err != nil {
			console.Info(err)
		} else {
			console.Info("kill process", processes[i].Pid, "success")
		}
	}
}
