/**
* @program: find-process
*
* @description:
*
* @author: lemo
*
* @create: 2022-09-08 14:53
**/

package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/lemonyxk/console"
	"github.com/lemonyxk/utils/v3"
)

var processes []*P

var flags = []string{
	"-a", "-c",
	"-o", "--port",
	"-p", "--pid",
	"-h", "--help",
	"-k", "--kill",
}

func main() {

	initProc()

	initPortMap()

	var args = filterArgs()

	var ps Processes

	// help
	if hasArgs("-h", "--help") {
		console.Info(help())
		return
	}

	// list
	if len(args) == 0 {
		ps = list()
		console.Info(ps)
		console.Info(fmt.Sprintf("[%d results]", len(ps)))
		return
	}

	switch argsIndex(1) {
	case "-o", "--port":
		ps = findProcessByPort(toInt32(os.Args[2:])...)
	// find process by pid
	case "-p", "--pid":
		ps = findProcessByPID(toInt32(os.Args[2:])...)
	default:
		ps = findProcessByString(filterArgs()...)
	}

	console.Info(ps)

	console.Info(fmt.Sprintf("[%d results]", len(ps)))

	if len(ps) == 0 {
		return
	}

	if !hasArgs("-k", "--kill") {
		return
	}

	kill(ps)
}

func filterArgs() []string {
	var res []string
	for i := 1; i < len(os.Args); i++ {
		if utils.ComparableArray(&flags).Has(os.Args[i]) {
			continue
		}
		res = append(res, os.Args[i])
	}
	return res
}

func toInt32(str []string) []int32 {
	if len(str) == 0 {
		return nil
	}

	var res []int32

	for i := 0; i < len(str); i++ {
		var i32, err = strconv.Atoi(str[i])
		if err != nil {
			continue
		}
		res = append(res, int32(i32))
	}

	return res
}

func argsIndex(index int) string {
	if len(os.Args) < index+1 {
		return ""
	}
	return os.Args[index]
}

func hasArgs(flag ...string) bool {
	var args = os.Args
	for i := 0; i < len(args); i++ {
		for j := 0; j < len(flag); j++ {
			if args[i] == flag[j] {
				return true
			}
		}
	}
	return false
}

func getArgs(flag ...string) string {
	var args = os.Args
	for i := 0; i < len(args); i++ {
		for j := 0; j < len(flag); j++ {
			if args[i] == flag[j] {
				if i+1 < len(args) {
					return args[i+1]
				}
			}
		}
	}
	return ""
}
