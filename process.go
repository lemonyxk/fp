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

	process   *P
	isColored bool
	colorStr  string
}

var pidProcess *process.Process

var selfPid = os.Getpid()

type Processes []Process

func (p Processes) String() string {

	if len(p) == 0 {
		return ""
	}

	var all = hasArgs("-a")
	var cmd = hasArgs("-c")

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

	if gidMaxLen > 0 {
		gidMaxLen += 2
	}
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
		if !all {
			if strings.HasPrefix(p[i].UserName, "_") {
				continue
			}
		}

		var add = ""

		// str += utils.Time.Timestamp(p[i].CreateTime).Format("01-02 15:04:05") +
		// 	strings.Repeat(" ", 2)

		add += p[i].Pid + strings.Repeat(" ", pidMaxLen-text.RuneCount(p[i].Pid))

		if p[i].GroupID != "" {
			add += p[i].GroupID + strings.Repeat(" ", gidMaxLen-text.RuneCount(p[i].GroupID))
		}

		add += p[i].Mem + strings.Repeat(" ", memMaxLen-text.RuneCount(p[i].Mem))

		add += p[i].UserName + strings.Repeat(" ", userMaxLen-text.RuneCount(p[i].UserName))

		if cmdMaxLen < 0 {
			nameMaxLen = termWidth - (gidMaxLen + pidMaxLen + portMaxLen + memMaxLen + userMaxLen)
			add += p[i].Name + strings.Repeat(" ", nameMaxLen-text.RuneCount(p[i].Name))
		} else {
			add += p[i].Name + strings.Repeat(" ", nameMaxLen-text.RuneCount(p[i].Name))
		}

		if p[i].Port != "" {
			add += p[i].Port + strings.Repeat(" ", portMaxLen-text.RuneCount(p[i].Port))
		}

		if p[i].Cmd != "" && cmdMaxLen > 0 {
			var cll = text.RuneCount(p[i].Cmd)
			if cmdMaxLen-cll > 0 {
				add += p[i].Cmd + strings.Repeat(" ", cmdMaxLen-cll)
			} else {
				if cmd {
					var index = 0
					var s []rune
					for {
						s = append(s, rune(p[i].Cmd[index]))
						if text.RuneCount(string(s)) >= cmdMaxLen {
							add += string(s) + "\n"
							add += strings.Repeat(" ", termWidth-cmdMaxLen)
							s = []rune{}
						}
						index++
						if index >= len(p[i].Cmd) {
							if len(s) != 0 {
								add += string(s)
							}
							break
						}
					}
				} else {
					var index = 0
					var l = len(p[i].Cmd)
					var s []rune
					for {
						s = append([]rune{rune(p[i].Cmd[l-index-1])}, s...)
						if text.RuneCount(string(s)) >= cmdMaxLen {
							add += string(s)
							break
						}
						index++
						if index >= len(p[i].Cmd) {
							if len(s) != 0 {
								add += string(s)
							}
							break
						}
					}
				}
			}
		}

		if i != len(p)-1 {
			add += "\n"
		}

		if p[i].isColored {
			add = strings.ReplaceAll(add, p[i].colorStr, console.FgRed.Sprint(p[i].colorStr))
		}

		str += add
	}

	return str

}

func list() Processes {
	var res []Process
	for i := 0; i < len(processes); i++ {
		var p = processes[i]
		var r, err = makeRes(p)
		if err != nil {
			continue
		}
		res = append(res, r)
	}

	return res
}

func findProcessByPID(pid ...int32) Processes {

	if len(pid) == 0 {
		return nil
	}

	var res []Process
	for i := 0; i < len(processes); i++ {
		var p = processes[i]
		if utils.ComparableArray(&pid).Has(p.Pid) {

			var r, err = makeRes(p)
			if err != nil {
				return nil
			}

			r.isColored = true
			r.colorStr = r.Pid

			res = append(res, r)

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
		var p = processes[i]

		var r, err = makeRes(p)
		if err != nil {
			continue
		}

		for j := 0; j < len(str); j++ {

			if strings.Contains(r.Name, str[j]) {
				r.isColored = true
				r.colorStr = str[j]
				res = append(res, r)
				break
			} else if strings.Contains(r.Pid, str[j]) {
				r.isColored = true
				r.colorStr = str[j]
				res = append(res, r)
				break
			} else if strings.Contains(r.Cmd, str[j]) {
				var pidS = getPid(pidProcess)
				if utils.ComparableArray(&pidS).Has(p.Pid) {
					break
				}
				r.isColored = true
				r.colorStr = str[j]
				res = append(res, r)
				break
			}
		}
	}

	return res
}

func makeRes(p *P) (Process, error) {
	var name, err = p.Name()
	if err != nil {
		return Process{}, err
	}

	un, _ := p.UserName()

	createTime, _ := p.CreateTime()

	cmd, _ := p.CmdLine()
	var mStr = ""
	mem, err := p.MemoryInfo()
	if err != nil {
		mStr = "deny"
	} else {
		mStr = size(int64(mem.RSS))
	}

	var gid = getGroupID(p)
	var gidStr = fmt.Sprintf("%d", gid)
	if gid == -1 {
		gidStr = ""
	}

	var rp = Process{
		Name:       name,
		Pid:        fmt.Sprintf("%d", p.Pid),
		CreateTime: createTime / 1000,
		process:    p,
		Cmd:        cmd,
		Mem:        mStr,
		UserName:   un,
		GroupID:    gidStr,
	}

	return rp, nil
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

	var pss Processes

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

		pss = append(pss, p...)

	}

	return pss
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
