// Copyright 2015 CodeIgnition. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

// +build linux

// Package netstat provides the network statistics
// of the system.
//
// On linux, it gets the data by executing the
// command `netstat -anp`. However, `netstat` is
// deprecated. So, If any problems arise, refactor
// the code to use `ss`.
package netstat

import (
	"os/exec"
	"strings"
)

// Data represents the network statistics data.
type Data []map[string]string

// CollectData collects the data and returns
// an error if any.
func CollectData() (Data, error) {
	var d Data
	out, err := exec.Command("netstat", "-anp").Output()
	if err != nil {
		return d, err
	}
	lines := strings.Split(string(out), "\n")

	// The output of `netstat -anp` has 2 sections,
	// Active Internet connections and Active UNIX domain sockets.
	// Both of them have different columns, so we need to know
	// which section does the line we are processing belongs to.
	var connType string

	for _, line := range lines {
		// Section heading
		if strings.HasPrefix(line, "Active Internet") {
			connType = "internet"
			continue
		}

		// Section heading
		if strings.HasPrefix(line, "Active UNIX") {
			connType = "unix"
			continue
		}

		// Column heading
		if strings.HasPrefix(line, "Proto") {
			continue
		}

		if connType == "internet" {
			internetConn(&d, line)
		}

	}
	return d, nil
}

// internetConn handles the lines under the section "Active Internet Connections"
// in the output of `netstat -anp`.
// We use *Data as we are appending in the function, thereby changing its internal
// length.
func internetConn(d *Data, line string) {
	// Proto Recv-Q Send-Q Local Address           Foreign Address         State       PID/Program name
	// tcp        0      0 127.0.1.1:53            0.0.0.0:*               LISTEN      -
	// udp        0      0 0.0.0.0:631             0.0.0.0:*                           -
	// udp        0      0 0.0.0.0:5353            0.0.0.0:*                           2251/chrome

	a := strings.Fields(line)

	var (
		state, pid, progname string
	)
	// return the length is less than 6,
	// to avoid a runtime panic when we access indices > 5
	if len(a) < 6 {
		return
	}

	// when state is empty
	if len(a) == 6 {
		state = ""
		// we ignore the case when a[5] is - , as pid and progname are
		// empty strings (zero values during declaration) already.
		if strings.Contains(a[5], "/") {
			b := strings.SplitN(a[5], "/", 2)
			pid = b[0]
			progname = b[1]
		}
	}

	if len(a) == 7 {
		state = a[5]
		if strings.Contains(a[6], "/") {
			b := strings.SplitN(a[6], "/", 2)
			pid = b[0]
			progname = b[1]
		}
	}

	m := map[string]string{
		"protocol":        a[0],
		"local_address":   a[3],
		"foreign_address": a[4],
		"state":           state,
		"process_id":      pid,
		"program_name":    progname,
	}
	*d = append(*d, m)
}
