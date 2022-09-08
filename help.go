/**
* @program: fp
*
* @description:
*
* @author: lemo
*
* @create: 2022-09-08 18:08
**/

package main

func help() string {
	return `
Usage: fp -o port1 port2
  -- find process by port number

Usage: fp -l
  -- list all processes

Usage: fp -p pid1 pid2
  -- find process by pid number

Usage: fp str1 str2
  -- find process by string
`
}
