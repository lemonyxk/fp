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
	"os"
	"strconv"

	"github.com/lemonyxk/console"
)

var processes []*P

func main() {

	initProc()

	initPortMap()

	if len(os.Args) == 1 {
		console.Info(list())
		return
	}

	var ps Processes

	switch os.Args[1] {
	// list all processes
	case "-a":
		console.Info(list())
		return
	// find process by port
	case "-o", "--port":
		ps = findProcessByPort(toInt32(os.Args[2:])...)
	// find process by pid
	case "-p", "--pid":
		ps = findProcessByPID(toInt32(os.Args[2:])...)
	case "-h", "--help":
		console.Info(help())
		return
	default:
		ps = findProcessByString(filterArgs()...)
	}

	console.Info(ps)

	if len(ps) == 0 {
		return
	}

	var k = hasArgs("-k")
	if !k {
		k = hasArgs("--kill")
	}

	if !k {
		return
	}

	kill(ps)
}

func filterArgs() []string {
	var index = -1
	for i := 0; i < len(os.Args); i++ {
		if os.Args[i] == "-k" || os.Args[i] == "--kill" {
			index = i
			break
		}
	}

	var args = os.Args

	if index != -1 {
		args = args[0:index]
	}
	return args[1:]
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

func hasArgs(flag string) bool {
	var args = os.Args
	for i := 0; i < len(args); i++ {
		if args[i] == flag {
			return true
		}
	}
	return false
}

func getArgs(flag []string) string {
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
