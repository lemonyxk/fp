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

	"github.com/lemonyxk/console"
	"github.com/lemonyxk/utils/v3"
)

type Process struct {
	Name       string
	Pid        int32
	CreateTime int64
}

type Processes []Process

func (p Processes) String() string {
	if len(p) == 0 {
		return "No results"
	}

	var table = console.NewTable()
	// table.Style().Options.DrawBorder = false
	// table.Style().Options.SeparateColumns = false
	// table.Header("CreateTime", "Pid", "Name")
	for i := 0; i < len(p); i++ {
		table.Row(
			utils.Time.Timestamp(p[i].CreateTime).Format("01-02 15:04:05"),
			p[i].Pid,
			p[i].Name,
		)
	}
	return table.Render()
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
			Pid:        process.Pid,
			CreateTime: createTime / 1000,
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

			res = append(res, Process{
				Name:       name,
				Pid:        process.Pid,
				CreateTime: createTime / 1000,
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
			Pid:        process.Pid,
			CreateTime: createTime / 1000,
		}

		for j := 0; j < len(str); j++ {
			if strings.Contains(name, str[j]) ||
				strings.Contains(fmt.Sprintf("%d", process.Pid), str[j]) {
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
