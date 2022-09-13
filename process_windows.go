//go:build windows
// +build windows

/**
* @program: find-process
*
* @description:
*
* @author: lemo
*
* @create: 2022-09-08 17:41
**/

package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/lemonyxk/console"
	"github.com/shirou/gopsutil/v3/process"
)

type P struct {
	*process.Process
	name string
}

func getPid(p *process.Process) []int32 {
	return []int{int32(os.Getpid())}
}

func (p *P) Name() (string, error) {
	return p.name, nil
}

func initProc() {
	console.SetFlags(0)
	console.Colorful(false)

	ps, err := process.Processes()
	if err != nil {
		console.Exit(err)
	}

	var m = getMap()

	for i := 0; i < len(ps); i++ {
		processes = append(processes, &P{Process: ps[i], name: m[ps[i].Pid]})
		if ps[i].Pid == int32(selfPid) {
			pidProcess = ps[i]
		}
	}

}

func getMap() map[int32]string {
	var re = regexp.MustCompile(`\s+`)
	var bts, err = execCmd("tasklist")
	if err != nil {
		return nil
	}

	var res = make(map[int32]string)

	var arr = strings.Split(string(bts), "\n")
	for i := 0; i < len(arr); i++ {

		var str = arr[i]
		str = strings.TrimLeft(str, " ")
		str = strings.TrimRight(str, " ")
		var findArr = re.Split(str, -1)

		if len(findArr) < 5 {
			continue
		}

		var pid, _ = strconv.Atoi(findArr[1])
		res[int32(pid)] = findArr[0]

	}

	return res
}

func findProcessByPort(port ...int32) Processes {
	if len(port) == 0 {
		return nil
	}

	var re = regexp.MustCompile(`\s+`)
	var bts, err = execCmd("netstat", "-navo")
	if err != nil {
		return nil
	}

	var processes Processes

	var arr = strings.Split(string(bts), "\n")
	for i := 0; i < len(arr); i++ {

		if !strings.Contains(arr[i], "LISTEN") {
			continue
		}

		var str = arr[i]
		str = strings.TrimLeft(str, " ")
		str = strings.TrimRight(str, " ")
		var findArr = re.Split(str, -1)

		if len(findArr) < 5 {
			continue
		}

		var ok = false
		var cp int32 = -1
		for j := 0; j < len(port); j++ {
			if strings.HasSuffix(findArr[1], fmt.Sprintf(":%d", port[j])) {
				ok = true
				cp = port[j]
				break
			}
		}

		if !ok {
			continue
		}

		intP, err := strconv.Atoi(findArr[4])
		if err != nil {
			continue
		}

		var p = findProcessByPID(int32(intP))
		if len(p) == 0 {
			continue
		}

		for k := 0; k < len(p); k++ {
			p[k].Port = console.FgRed.Sprintf("%d", cp)
			p[k].Pid = fmt.Sprintf("%d", p[k].process.Pid)
		}

		processes = append(processes, p...)

	}

	return processes
}
