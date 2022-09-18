//go:build linux
// +build linux

/**
* @program: fp
*
* @description:
*
* @author: lemo
*
* @create: 2022-09-15 14:58
**/

package main

import (
	"strconv"
	"strings"
	"syscall"

	"github.com/lemonyxk/console"
)

var netMap = make(map[int]int)

func initPortMap() {
	var str, err = execCmd("netstat", "-nap")
	if err != nil {
		console.Exit(err)
	}

	var res = getArrFromLineStr(string(str), []string{"LISTEN", "tcp"}, nil)
	for i := 0; i < len(res); i++ {
		var addr = res[i][3]
		var port = addr[strings.LastIndex(addr, ":")+1:]

		var o, _ = strconv.Atoi(port)
		var p, _ = strconv.Atoi(res[i][6][:strings.Index(res[i][6], "/")])
		netMap[o] = p
	}
}

func getGroupID(p *P) int {
	var g, err = syscall.Getpgid(int(p.Pid))
	if err != nil {
		return syscall.Getppid()
	}
	return g
}
