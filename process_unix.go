//go:build linux || unix || darwin
// +build linux unix darwin

/**
* @program: find-process
*
* @description:
*
* @author: lemo
*
* @create: 2022-09-08 17:42
**/

package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/lemonyxk/console"
	"github.com/shirou/gopsutil/v3/process"
)

type P struct {
	*process.Process
}

func (p *P) Name() (string, error) {
	return p.Process.Name()
}

func getPid(p *process.Process) []int32 {

	var res []int32

	res = append(res, p.Pid)

	for {
		pp, err := p.Parent()
		if err != nil {
			break
		}

		n, err := p.Name()
		if err != nil {
			break
		}

		if strings.ToUpper(n) == "SUDO" {
			break
		}

		res = append(res, pp.Pid)

		p = pp
	}

	return res
}

func initProc() {
	console.SetFlags(0)
	console.Colorful(false)

	ps, err := process.Processes()
	if err != nil {
		console.Exit(err)
	}

	for i := 0; i < len(ps); i++ {
		processes = append(processes, &P{ps[i]})
		if ps[i].Pid == int32(selfPid) {
			pidProcess = ps[i]
		}
	}
}

func findProcessByPort(port ...int32) Processes {
	if len(port) == 0 {
		return nil
	}

	var processes Processes

	for i := 0; i < len(port); i++ {

		var str, err = execCmd("lsof", "-t", "-s", "TCP:LISTEN", "-i", fmt.Sprintf(":%d", port[i]))
		if err != nil {
			continue
		}

		var arr = strings.Split(string(str), "\n")
		for j := 0; j < len(arr); j++ {
			var s = strings.TrimSpace(arr[j])
			if s == "" {
				continue
			}

			var intP, _ = strconv.Atoi(s)
			if intP == 0 {
				continue
			}

			var p = findProcessByPID(int32(intP))
			if len(p) == 0 {
				continue
			}

			for k := 0; k < len(p); k++ {
				p[k].Port = console.FgRed.Sprintf("%d", port[i])
				p[k].Pid = fmt.Sprintf("%d", p[k].process.Pid)
			}

			processes = append(processes, p...)
		}

	}

	return processes
}
