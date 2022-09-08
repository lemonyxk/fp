//go:build unix
// +build unix

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
)

func findProcessByPort(port ...int32) Processes {
	if len(port) == 0 {
		return nil
	}

	var ps []int32

	var params = []string{"lsof"}
	for i := 0; i < len(port); i++ {
		params = append(params, "-i", fmt.Sprintf(":%d", port[i]))
	}

	params = append(params, "-t", "-s", "TCP:LISTEN")

	var str, err = execCmd(params[0], params[1:]...)
	if err != nil {
		return nil
	}

	var arr = strings.Split(string(str), "\n")
	for i := 0; i < len(arr); i++ {
		var s = strings.TrimSpace(arr[i])
		if s == "" {
			continue
		}

		var intP, _ = strconv.Atoi(s)
		if intP == 0 {
			continue
		}

		ps = append(ps, int32(intP))
	}

	return findProcessByPID(ps...)
}
