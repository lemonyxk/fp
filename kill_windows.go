//go:build windows
// +build windows

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
	"github.com/lemonyxk/console"
)

func kill(processes Processes) {
	for i := 0; i < len(processes); i++ {
		var err = processes[i].process.Kill()
		if err != nil {
			console.Info(err)
		} else {
			console.Info("kill process", processes[i].Pid, "success")
		}
	}
}
