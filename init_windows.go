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
	"os"
	"path/filepath"
	"strings"

	"github.com/lemonyxk/console"
	"github.com/shirou/gopsutil/v3/process"
)

type P struct {
	*process.Process
	// name string
}

func getPid(p *process.Process) []int32 {
	return []int32{int32(os.Getpid())}
}

func (p *P) Name() (string, error) {
	var exe, err = p.Exe()
	return filepath.Base(exe), err
}

func (p *P) UserName() (string, error) {
	var name, _ = p.Process.Username()
	var s = filepath.Base(name)
	if s == "." {
		return "deny", nil
	}
	return s, nil
}

func (p *P) CmdLine() (string, error) {
	var cmd, err = p.Process.Cmdline()
	return strings.ReplaceAll(cmd, `"`, ""), err
}

func initProc() {
	console.SetFlags(0)
	console.Colorful(false)

	ps, err := process.Processes()
	if err != nil {
		console.Exit(err)
	}

	// var m = getMap()

	for i := 0; i < len(ps); i++ {
		// processes = append(processes, &P{Process: ps[i], name: m[ps[i].Pid]})
		processes = append(processes, &P{Process: ps[i]})
		if ps[i].Pid == int32(selfPid) {
			pidProcess = ps[i]
		}
	}

}

// func getMap() map[int32]string {
// 	var re = regexp.MustCompile(`\s+`)
// 	var bts, err = execCmd("tasklist")
// 	if err != nil {
// 		return nil
// 	}
//
// 	var res = make(map[int32]string)
//
// 	var arr = strings.Split(string(bts), "\n")
// 	for i := 0; i < len(arr); i++ {
//
// 		var str = arr[i]
// 		str = strings.TrimLeft(str, " ")
// 		str = strings.TrimRight(str, " ")
// 		var findArr = re.Split(str, -1)
//
// 		if len(findArr) < 5 {
// 			continue
// 		}
//
// 		var pid, _ = strconv.Atoi(findArr[1])
// 		res[int32(pid)] = findArr[0]
//
// 	}
//
// 	return res
// }
