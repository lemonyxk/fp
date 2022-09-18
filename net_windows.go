//go:build windows
// +build windows

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

	"github.com/lemonyxk/console"
)

var netMap = make(map[int]int)

func initPortMap() {
	var str, err = execCmd("netstat", "-anvo")
	if err != nil {
		console.Exit(err)
	}

	var res = getArrFromLineStr(string(str), []string{"LISTEN", "TCP"}, nil)
	for i := 0; i < len(res); i++ {
		var addr = res[i][1]
		var port = addr[strings.LastIndex(addr, ":")+1:]

		var o, _ = strconv.Atoi(port)
		var p, _ = strconv.Atoi(res[i][4])
		netMap[o] = p
	}
}

func getGroupID(p *P) int {
	// var r, _ = p.Ppid()
	// return int(r)
	return -1
}
