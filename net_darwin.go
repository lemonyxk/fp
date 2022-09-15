//go:build darwin
// +build darwin

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
	var str, err = execCmd("netstat", "-anv")
	if err != nil {
		console.Exit(err)
	}

	var res = getArrFromLineStr(string(str), []string{"LISTEN"}, nil)
	for i := 0; i < len(res); i++ {
		var addr = res[i][3]
		var port = addr[strings.LastIndex(addr, ".")+1:]

		var o, _ = strconv.Atoi(port)
		var p, _ = strconv.Atoi(res[i][8])
		netMap[o] = p
	}
}
