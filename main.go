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

	switch os.Args[1] {
	// list all processes
	case "-l":
		console.Info(list())
	// find process by port
	case "-i":
		console.Info(findProcessByPort(toInt32(os.Args[2:])...))
	// find process by pid
	case "-p":
		console.Info(findProcessByPID(toInt32(os.Args[2:])...))
	default:
		console.Info(findProcessByString(os.Args[1:]...))
	}
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
