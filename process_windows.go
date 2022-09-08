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
)

func findProcessByPort(port ...int32) Processes {
	if len(port) == 0 {
		return nil
	}

	var ps []int32

	var re = regexp.MustCompile(`\s+`)
	var bts, err = execCmd("netstat", "-navo")
	if err != nil {
		return nil
	}

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
		for j := 0; j < len(port); j++ {
			if strings.HasSuffix(findArr[1], fmt.Sprintf(":%d", port[j])) {
				ok = true
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

		ps = append(ps, int32(intP))
	}

	return findProcessByPID(ps...)

}
