/**
* @program: find-process
*
* @description:
*
* @author: lemo
*
* @create: 2022-09-08 15:02
**/

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/text"
	"github.com/lemonyxk/console"
	"github.com/lemonyxk/utils/v3"
	"github.com/olekukonko/ts"
)

type Process struct {
	Name       string
	Pid        string
	Port       string
	CreateTime int64
	Cmd        string
	Mem        string
	UserName   string
	process    *P
}

type Processes []Process

func (p Processes) String() string {

	if len(p) == 0 {
		return "No results"
	}

	var size, err = ts.GetSize()
	if err != nil {
		console.Exit(err)
	}

	var termWidth = size.Col()

	var now = time.Now().Format("01-02 15:04:05")

	var timeMaxLen = text.RuneCount(now)
	var pidMaxLen = 0
	var userMaxLen = 0
	var nameMaxLen = 0
	var portMaxLen = 0
	var memMaxLen = 0
	var cmdMaxLen = 0

	for i := 0; i < len(p); i++ {
		var px = text.RuneCount(p[i].Pid)
		if px > pidMaxLen {
			pidMaxLen = px
		}

		var ux = text.RuneCount(p[i].UserName)
		if ux > userMaxLen {
			userMaxLen = ux
		}

		var nx = text.RuneCount(p[i].Name)
		if nx > nameMaxLen {
			nameMaxLen = nx
		}

		var ox = text.RuneCount(p[i].Port)
		if ox > portMaxLen {
			portMaxLen = ox
		}

		var mx = text.RuneCount(p[i].Mem)
		if mx > memMaxLen {
			memMaxLen = mx
		}

		var cx = text.RuneCount(p[i].Cmd)
		if cx > cmdMaxLen {
			cmdMaxLen = cx
		}
	}

	timeMaxLen += 4
	pidMaxLen += 4
	nameMaxLen += 4
	portMaxLen += 4
	memMaxLen += 4
	userMaxLen += 4

	if timeMaxLen+pidMaxLen+nameMaxLen+portMaxLen+memMaxLen+cmdMaxLen+userMaxLen > termWidth {
		cmdMaxLen = termWidth - (timeMaxLen + pidMaxLen + nameMaxLen + portMaxLen + memMaxLen + userMaxLen)
	}

	var str = ""

	for i := 0; i < len(p); i++ {
		str += utils.Time.Timestamp(p[i].CreateTime).Format("01-02 15:04:05") +
			strings.Repeat(" ", 4)

		str += p[i].UserName + strings.Repeat(" ", userMaxLen-text.RuneCount(p[i].UserName))

		str += p[i].Pid + strings.Repeat(" ", pidMaxLen-text.RuneCount(p[i].Pid))

		str += p[i].Name + strings.Repeat(" ", nameMaxLen-text.RuneCount(p[i].Name))

		str += p[i].Mem + strings.Repeat(" ", memMaxLen-text.RuneCount(p[i].Mem))

		if p[i].Port != "" {
			str += p[i].Port + strings.Repeat(" ", portMaxLen-text.RuneCount(p[i].Port))
		}

		if p[i].Cmd != "" {
			if cmdMaxLen-text.RuneCount(p[i].Cmd) > 0 {
				str += p[i].Cmd + strings.Repeat(" ", cmdMaxLen-text.RuneCount(p[i].Cmd))
			} else {
				str += p[i].Cmd[:cmdMaxLen-3] + "..."
			}
		}

		str += "\n"
	}

	return str

}

func list() Processes {
	var res []Process
	for i := 0; i < len(processes); i++ {
		var process = processes[i]
		var name, err = process.Name()
		if err != nil {
			continue
		}

		createTime, err := process.CreateTime()
		if err != nil {
			continue
		}

		res = append(res, Process{
			Name:       name,
			Pid:        fmt.Sprintf("%d", process.Pid),
			CreateTime: createTime / 1000,
			process:    process,
		})
	}

	return res
}

func findProcessByPID(pid ...int32) Processes {

	if len(pid) == 0 {
		return nil
	}

	var res []Process
	for i := 0; i < len(processes); i++ {
		var process = processes[i]
		if utils.ComparableArray(&pid).Has(process.Pid) {
			var name, err = process.Name()
			if err != nil {
				return nil
			}

			createTime, err := process.CreateTime()
			if err != nil {
				return nil
			}

			cmd, _ := process.Cmdline()
			var mStr = ""
			mem, err := process.MemoryInfo()
			if err != nil {
				mStr = "deny"
			} else {
				mStr = size(int64(mem.RSS))
			}

			un, _ := process.Username()

			res = append(res, Process{
				Name:       name,
				Pid:        console.FgRed.Sprintf("%d", process.Pid),
				CreateTime: createTime / 1000,
				process:    process,
				Cmd:        cmd,
				Mem:        mStr,
				UserName:   un,
			})

			if len(pid) == len(res) {
				return res
			}
		}
	}

	return res
}

func findProcessByString(str ...string) Processes {

	if len(str) == 0 {
		return nil
	}

	var res []Process
	for i := 0; i < len(processes); i++ {
		var process = processes[i]
		var name, err = process.Name()
		if err != nil {
			continue
		}

		createTime, err := process.CreateTime()
		if err != nil {
			continue
		}

		var r = Process{
			Name:       name,
			Pid:        fmt.Sprintf("%d", process.Pid),
			CreateTime: createTime / 1000,
			process:    process,
		}

		for j := 0; j < len(str); j++ {
			if strings.Contains(name, str[j]) {
				r.Name = strings.Replace(name, str[j], console.FgRed.Sprintf("%s", str[j]), 1)
				cmd, _ := process.Cmdline()
				un, _ := process.Username()

				var mStr = ""
				mem, err := process.MemoryInfo()
				if err != nil {
					mStr = "deny"
				} else {
					mStr = size(int64(mem.RSS))
				}

				r.Cmd = cmd
				r.Mem = mStr
				r.UserName = un
				res = append(res, r)
				break
			} else if strings.Contains(fmt.Sprintf("%d", process.Pid), str[j]) {
				r.Pid = strings.Replace(fmt.Sprintf("%d", process.Pid), str[j],
					console.FgRed.Sprintf("%s", str[j]), 1)
				cmd, _ := process.Cmdline()
				un, _ := process.Username()

				var mStr = ""
				mem, err := process.MemoryInfo()
				if err != nil {
					mStr = "deny"
				} else {
					mStr = size(int64(mem.RSS))
				}

				r.Cmd = cmd
				r.Mem = mStr
				r.UserName = un
				res = append(res, r)
				break
			}
		}
	}

	return res
}

func execCmd(c string, args ...string) ([]byte, error) {

	var buf bytes.Buffer

	cmd := exec.Command(c, args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = &buf
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func size(i int64) string {

	var s = float64(i)

	if s < 1024*1024 {
		return fmt.Sprintf("%.1fKB", s/1024)
	}

	return fmt.Sprintf("%.1fMB", s/1024/1024)
}
