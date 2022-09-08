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
	"github.com/shirou/gopsutil/v3/process"
)

var processes []*process.Process

func init() {
	console.SetFlags(0)
	console.Colorful(false)

	var err error
	processes, err = process.Processes()
	if err != nil {
		console.Exit(err)
	}
}

func main() {

	if len(os.Args) < 2 {
		console.Exit(help())
	}

	var processes Processes

	switch os.Args[1] {
	// list all processes
	case "-l", "--list":
		console.Info(list())
	// find process by port
	case "-o", "--port":
		processes = findProcessByPort(toInt32(os.Args[2:])...)
	// find process by pid
	case "-p", "--pid":
		processes = findProcessByPID(toInt32(os.Args[2:])...)
	default:
		processes = findProcessByString(filterArgs()...)
	}

	console.Info(processes)

	if len(processes) == 0 {
		return
	}

	var k = hasArgs("-k", os.Args)
	if !k {
		k = hasArgs("--kill", os.Args)
	}

	if !k {
		return
	}

	kill(processes)
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

func hasArgs(flag string, args []string) bool {
	for i := 0; i < len(args); i++ {
		if args[i] == flag {
			return true
		}
	}
	return false
}

func getArgs(flag []string, args []string) string {
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
