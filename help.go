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
Usage: fp -o port
  -- find process by port number

Usage: fp
  -- current and root user processes

Usage: fp -a
  -- all processes

Usage: fp -c
  -- full command info

Usage: fp -p pid
  -- find process by pid number

Usage: fp string
  -- find process by string`
}
