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
	"regexp"
	"strings"

	"github.com/jedib0t/go-pretty/text"
	"github.com/lemonyxk/console"
	"github.com/lemonyxk/utils/v3"
	"github.com/olekukonko/ts"
	"github.com/shirou/gopsutil/v3/process"
)

type Process struct {
	Name       string
	Pid        string
	Port       string
	CreateTime int64
	Cmd        string
	Mem        string
	UserName   string
	GroupID    string
	process    *P
}

var pidProcess *process.Process

var selfPid = os.Getpid()

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

	// var now = time.Now().Format("01-02 15:04:05")

	// var timeMaxLen = text.RuneCount(now)

	var gidMaxLen = 0
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

		var gx = text.RuneCount(p[i].GroupID)
		if gx > gidMaxLen {
			gidMaxLen = gx
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

	gidMaxLen += 2
	pidMaxLen += 2
	nameMaxLen += 2
	if portMaxLen > 0 {
		portMaxLen += 2
	}
	memMaxLen += 2
	userMaxLen += 2

	cmdMaxLen = termWidth - (gidMaxLen + pidMaxLen + nameMaxLen + portMaxLen + memMaxLen + userMaxLen)

	var str = ""

	for i := 0; i < len(p); i++ {
		// str += utils.Time.Timestamp(p[i].CreateTime).Format("01-02 15:04:05") +
		// 	strings.Repeat(" ", 2)

		str += p[i].Pid + strings.Repeat(" ", pidMaxLen-text.RuneCount(p[i].Pid))

		str += p[i].GroupID + strings.Repeat(" ", gidMaxLen-text.RuneCount(p[i].GroupID))

		str += p[i].Mem + strings.Repeat(" ", memMaxLen-text.RuneCount(p[i].Mem))

		str += p[i].UserName + strings.Repeat(" ", userMaxLen-text.RuneCount(p[i].UserName))

		if cmdMaxLen < 0 {
			nameMaxLen = termWidth - (gidMaxLen + pidMaxLen + portMaxLen + memMaxLen + userMaxLen)
			str += p[i].Name + strings.Repeat(" ", nameMaxLen-text.RuneCount(p[i].Name))
		} else {
			str += p[i].Name + strings.Repeat(" ", nameMaxLen-text.RuneCount(p[i].Name))
		}

		if p[i].Port != "" {
			str += p[i].Port + strings.Repeat(" ", portMaxLen-text.RuneCount(p[i].Port))
		}

		if p[i].Cmd != "" && cmdMaxLen > 0 {
			var cll = text.RuneCount(p[i].Cmd)
			if cmdMaxLen-cll > 0 {
				str += p[i].Cmd + strings.Repeat(" ", cmdMaxLen-cll)
			} else {
				str += "..." + p[i].Cmd[cll-cmdMaxLen+4:]
			}
		}

		if i != len(p)-1 {
			str += "\n"
		}
	}

	return str

}

func list() Processes {
	var res []Process
	var all = false
	if hasArgs("-a") {
		all = true
	}
	for i := 0; i < len(processes); i++ {
		var process = processes[i]
		var name, err = process.Name()
		if err != nil {
			continue
		}

		un, _ := process.Username()
		un = shortName(un)

		if !all {
			if strings.HasPrefix(un, "_") {
				continue
			}
		}

		createTime, _ := process.CreateTime()

		cmd, _ := process.Cmdline()
		var mStr = ""
		mem, err := process.MemoryInfo()
		if err != nil {
			mStr = "deny"
		} else {
			mStr = size(int64(mem.RSS))
		}

		res = append(res, Process{
			Name:       name,
			Pid:        fmt.Sprintf("%d", process.Pid),
			CreateTime: createTime / 1000,
			process:    process,
			Cmd:        cmd,
			Mem:        mStr,
			UserName:   un,
			GroupID:    fmt.Sprintf("%d", getGroupID(process)),
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

			createTime, _ := process.CreateTime()

			cmd, _ := process.Cmdline()
			var mStr = ""
			mem, err := process.MemoryInfo()
			if err != nil {
				mStr = "deny"
			} else {
				mStr = size(int64(mem.RSS))
			}

			un, _ := process.Username()
			un = shortName(un)

			res = append(res, Process{
				Name:       name,
				Pid:        console.FgRed.Sprintf("%d", process.Pid),
				CreateTime: createTime / 1000,
				process:    process,
				Cmd:        cmd,
				Mem:        mStr,
				UserName:   un,
				GroupID:    fmt.Sprintf("%d", getGroupID(process)),
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

		createTime, _ := process.CreateTime()

		var r = Process{
			Name:       name,
			Pid:        fmt.Sprintf("%d", process.Pid),
			CreateTime: createTime / 1000,
			process:    process,
		}

		for j := 0; j < len(str); j++ {

			cmd, _ := process.Cmdline()

			if strings.Contains(name, str[j]) {
				r.Name = strings.Replace(name, str[j], console.FgRed.Sprintf("%s", str[j]), 1)
				un, _ := process.Username()
				un = shortName(un)

				var mStr = ""
				mem, err := process.MemoryInfo()
				if err != nil {
					mStr = "deny"
				} else {
					mStr = size(int64(mem.RSS))
				}
				r.GroupID = fmt.Sprintf("%d", getGroupID(process))
				r.Cmd = cmd
				r.Mem = mStr
				r.UserName = un
				res = append(res, r)
				break
			} else if strings.Contains(fmt.Sprintf("%d", process.Pid), str[j]) {
				r.Pid = strings.Replace(fmt.Sprintf("%d", process.Pid), str[j],
					console.FgRed.Sprintf("%s", str[j]), 1)
				un, _ := process.Username()
				un = shortName(un)

				var mStr = ""
				mem, err := process.MemoryInfo()
				if err != nil {
					mStr = "deny"
				} else {
					mStr = size(int64(mem.RSS))
				}
				r.GroupID = fmt.Sprintf("%d", getGroupID(process))
				r.Cmd = cmd
				r.Mem = mStr
				r.UserName = un
				res = append(res, r)
				break
			} else if strings.Contains(cmd, str[j]) {

				var pids = getPid(pidProcess)

				if utils.ComparableArray(&pids).Has(process.Pid) {
					break
				}

				r.Cmd = cmd
				un, _ := process.Username()
				un = shortName(un)

				var mStr = ""
				mem, err := process.MemoryInfo()
				if err != nil {
					mStr = "deny"
				} else {
					mStr = size(int64(mem.RSS))
				}

				r.GroupID = fmt.Sprintf("%d", getGroupID(process))
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

func findProcessByPort(port ...int32) Processes {
	if len(port) == 0 {
		return nil
	}

	var processes Processes

	for i := 0; i < len(port); i++ {

		var intP, ok = netMap[int(port[i])]
		if !ok {
			continue
		}

		var p = findProcessByPID(int32(intP))
		if len(p) == 0 {
			continue
		}

		for k := 0; k < len(p); k++ {
			p[k].Port = console.FgRed.Sprintf("%d", port[i])
			p[k].Pid = fmt.Sprintf("%d", p[k].process.Pid)
		}

		processes = append(processes, p...)

	}

	return processes
}

func size(i int64) string {

	var s = float64(i)

	if s < 1024*1024 {
		return fmt.Sprintf("%.1fKB", s/1024)
	}

	return fmt.Sprintf("%.1fMB", s/1024/1024)
}

var re = regexp.MustCompile(`\s+`)

func getArrFromLineStr(str string, find, filter []string) [][]string {
	var res [][]string
	var arr = strings.Split(str, "\n")
	for i := 0; i < len(arr); i++ {

		var s = arr[i]
		s = strings.TrimLeft(s, " ")
		s = strings.TrimRight(s, " ")

		var fi = true
		for j := 0; j < len(filter); j++ {
			if strings.Contains(s, filter[j]) {
				fi = false
				break
			}
		}

		if !fi {
			continue
		}

		if len(find) == 0 {
			res = append(res, re.Split(s, -1))
			break
		}

		var f = true
		for j := 0; j < len(find); j++ {
			if !strings.Contains(s, find[j]) {
				f = false
				break
			}
		}

		if f {
			res = append(res, re.Split(s, -1))
		}
	}

	return res
}
