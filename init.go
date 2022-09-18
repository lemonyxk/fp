//go:build !windows
// +build !windows

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

func (p *P) UserName() (string, error) {
	return p.Process.Username()
}

func (p *P) CmdLine() (string, error) {
	return p.Process.Cmdline()
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
