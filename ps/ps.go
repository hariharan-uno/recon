// Copyright 2015 Hari haran. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

// +build linux

// Package ps provides data of the processes
// running on the system.
//
// On linux, it gets the data by executing the
// command `ps -aux`.
package ps

import (
	"os/exec"
	"strings"
)

// Data represents processes data.
type Data map[string][]map[string]string

// CollectData collects the data and returns
// an error if any.
func CollectData() (Data, error) {
	d := make(Data)
	out, err := exec.Command("ps", "aux").Output()
	if err != nil {
		return d, err
	}
	lines := strings.Split(string(out), "\n")

	// lines[0] is the column headings
	// USER       PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
	// root         1  0.0  0.0  33776  4264 ?        Ss   10:07   0:01 /sbin/init

	for _, line := range lines[1:] {
		a := strings.Fields(line)
		if len(a) >= 10 {
			// We namespace the processes by the user
			user := a[0]
			if _, ok := d[user]; !ok {
				// we don't know how many processes are going to be there for this user,
				// so we can't allocate the slice using make. The only option is to use
				// append.
				d[user] = []map[string]string{}
			}
			m := map[string]string{
				"process_id":                 a[1],
				"percentage_cpu_used":        a[2],
				"percentage_mem_used":        a[3],
				"virtual_memory_used":        a[4],
				"real_memory_used":           a[5],
				"terminal":                   a[6],
				"status_code":                a[7],
				"start_time":                 a[8],
				"total_cpu_utilization_time": a[9],
				"command":                    strings.Join(a[10:], " "),
			}
			d[user] = append(d[user], m)
		}
	}

	return d, nil
}
